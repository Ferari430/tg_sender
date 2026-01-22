package pkg

func RemoveDuplicates(input []string) []string {
	// Создаем map для отслеживания уникальных значений
	seen := make(map[string]bool)
	result := []string{}

	// Проходим по всем элементам
	for _, value := range input {
		// Если элемента еще нет в map, добавляем его в результат
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}

	return result
}
