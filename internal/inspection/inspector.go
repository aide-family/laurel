package inspection

type Inspector interface {
	Show()
}

func Show(inspectors ...Inspector) {
	for _, inspector := range inspectors {
		inspector.Show()
	}
}
