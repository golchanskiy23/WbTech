package postgres

type Option func(*Source)

func SetMaxPoolSize(size int) Option {
	return func(s *Source) {
		s.MaxPoolSize = size
	}
}
