package internal

type MIMEs []string

func (m MIMEs) Has(mime string) bool {
	for _, current := range m {
		if mime == current {
			return true
		}
	}
	return false
}
