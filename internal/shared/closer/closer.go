package closer

import (
	"github.com/1Storm3/flibox-api/internal/shared/logger"
	"os"
	"os/signal"
	"sync"

	"go.uber.org/zap"
)

// создаем глобальный closer
var globalCloser = New()

// Add добавляет функцию закрытия в глобальный closer
func Add(f ...func() error) {
	globalCloser.Add(f...)
}

// Wait ждет, пока все функции закрытия завершатся
func Wait() {
	globalCloser.Wait()
}

// CloseAll запускает все функции закрытия
func CloseAll() {
	globalCloser.CloseAll()
}

// Closer управляет списком функций закрытия
type Closer struct {
	mu    sync.Mutex
	once  sync.Once
	done  chan struct{}
	funcs []func() error
}

// New возвращает новый Closer и обрабатывает системные сигналы для вызова CloseAll
func New(sig ...os.Signal) *Closer {
	c := &Closer{done: make(chan struct{})}
	if len(sig) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, sig...)
			<-ch
			signal.Stop(ch)
			c.CloseAll()
		}()
	}
	return c
}

// Add добавляет одну или несколько функций закрытия в closer
func (c *Closer) Add(f ...func() error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.funcs = append(c.funcs, f...)
}

// Wait блокирует выполнение до завершения всех функций закрытия
func (c *Closer) Wait() {
	<-c.done
}

// CloseAll вызывает все функции закрытия
func (c *Closer) CloseAll() {
	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		var wg sync.WaitGroup
		errs := make(chan error, len(funcs))

		// Запускаем каждую функцию закрытия в отдельной горутине
		for _, f := range funcs {
			wg.Add(1)
			go func(f func() error) {
				defer wg.Done()
				if err := f(); err != nil {
					errs <- err
				}
			}(f)
		}

		// Закрываем канал errs, как только все горутины завершат выполнение
		go func() {
			wg.Wait()
			close(errs)
		}()

		// Читаем ошибки из канала и выводим их
		for err := range errs {
			if err != nil {
				logger.Info("Ошибка при закрытии:", zap.Error(err))
			}
		}
	})
}
