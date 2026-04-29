package ports

type Notifier interface {
	Notify(message string) error
}
