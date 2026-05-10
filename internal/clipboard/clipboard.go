package clipboard

func Copy(text string) error {
	return copyToClipboard(text)
}
