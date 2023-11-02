package models

func runModelValFuncs[T Pass](facet *T, fns ...func(*T) error) error {
	for _, fn := range fns {
		if err := fn(facet); err != nil {
			return err
		}
	}
	return nil
}
