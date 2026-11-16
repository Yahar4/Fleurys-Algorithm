package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Yahar4/internal/graph"
)

func main() {
	var g *graph.Graph

	if len(os.Args) > 1 {
		var err error
		g, err = loadFromFile(os.Args[1])
		if err != nil {
			fmt.Println("Ошибка:", err)
			os.Exit(1)
		}
		fmt.Println("Загружен граф из", os.Args[1])
	} else {
		g = readFromStdin()
	}

	analyze(g)
}

func analyze(g *graph.Graph) {
	edgesCount := 0
	for _, e := range g.GetEdges() {
		if e.IsActive() {
			edgesCount++
		}
	}
	fmt.Printf("\nВершин: %d, рёбер: %d\n", g.GetN(), edgesCount)
	for i := 0; i < g.GetN(); i++ {
		fmt.Printf("  deg(%d) = %d\n", i, g.Degree(i))
	}

	if !g.IsEulerian() {
		fmt.Println("\nГраф НЕ содержит эйлерова цикла")
		oddFound := false
		for i := 0; i < g.GetN(); i++ {
			if g.Degree(i)%2 != 0 {
				fmt.Printf("  причина: вершина %d имеет нечётную степень %d\n", i, g.Degree(i))
				oddFound = true
			}
		}
		if !oddFound {
			fmt.Println("  причина: граф несвязен")
		}
		return
	}

	fmt.Println("\nГраф содержит эйлеров цикл")
	start := 0
	for i := 0; i < g.GetN(); i++ {
		if g.Degree(i) > 0 {
			start = i
			break
		}
	}

	cycle := g.Fleury(start)

	fmt.Printf("Эйлеров цикл (старт из %d):\n", start)
	fmt.Println(" ", cycleString(cycle))
	fmt.Printf("Длина: %d вершин, %d рёбер (ожидалось %d рёбер)\n",
		len(cycle), len(cycle)-1, edgesCount)
}

func cycleString(cycle []int) string {
	parts := make([]string, len(cycle))
	for i, v := range cycle {
		parts[i] = strconv.Itoa(v)
	}
	return strings.Join(parts, " - ")
}

func loadFromFile(path string) (*graph.Graph, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var lines []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		lines = append(lines, line)
	}
	if len(lines) == 0 {
		return nil, fmt.Errorf("пустой файл")
	}

	var n, startRow int
	if !strings.Contains(lines[0], ":") {
		n, err = strconv.Atoi(lines[0])
		if err != nil || n < 0 {
			return nil, fmt.Errorf("некорректное n: %q", lines[0])
		}
		startRow = 1
	} else {
		n = 0
		for _, line := range lines {
			colon := strings.IndexByte(line, ':')
			if colon < 0 {
				continue
			}
			from, err := strconv.Atoi(strings.TrimSpace(line[:colon]))
			if err != nil {
				continue
			}
			if from+1 > n {
				n = from + 1
			}
		}
		startRow = 0
	}

	g := graph.New(n)
	seen := make(map[[2]int]bool)

	for _, line := range lines[startRow:] {
		colon := strings.IndexByte(line, ':')
		if colon < 0 {
			continue
		}
		from, err := strconv.Atoi(strings.TrimSpace(line[:colon]))
		if err != nil || from < 0 || from >= n {
			return nil, fmt.Errorf("некорректный номер вершины: %q", line[:colon])
		}

		rest := strings.TrimSpace(line[colon+1:])
		rest = strings.TrimPrefix(rest, "[")
		rest = strings.TrimSuffix(rest, "]")
		sep := func(r rune) bool { return r == ' ' || r == ',' || r == '\t' }

		for _, s := range strings.FieldsFunc(rest, sep) {
			to, err := strconv.Atoi(s)
			if err != nil || to < 0 || to >= n {
				return nil, fmt.Errorf("некорректный сосед %q у вершины %d", s, from)
			}
			if from == to {
				continue
			}
			a, b := from, to
			if a > b {
				a, b = b, a
			}
			key := [2]int{a, b}
			if seen[key] {
				continue
			}
			seen[key] = true
			g.AddEdge(from, to)
		}
	}
	return g, nil
}

func readFromStdin() *graph.Graph {
	r := bufio.NewReader(os.Stdin)

	fmt.Print("Количество вершин: ")
	var n int
	if _, err := fmt.Fscan(r, &n); err != nil || n < 0 {
		fmt.Println("Ошибка ввода n")
		os.Exit(1)
	}
	r.ReadString('\n')

	g := graph.New(n)
	seen := make(map[[2]int]bool)

	fmt.Printf("Введите %d строк вида \"i: сосед1 сосед2 ...\"\n", n)
	for i := 0; i < n; i++ {
		fmt.Printf("%d> ", i)
		line, _ := r.ReadString('\n')
		line = strings.TrimSpace(line)
		colon := strings.IndexByte(line, ':')
		if colon < 0 {
			fmt.Println("  ошибка: нет ':' — пропущено")
			continue
		}
		from, err := strconv.Atoi(strings.TrimSpace(line[:colon]))
		if err != nil || from < 0 || from >= n {
			fmt.Println("  ошибка: некорректный номер вершины")
			continue
		}
		rest := strings.TrimSuffix(strings.TrimPrefix(line[colon+1:], "["), "]")
		sep := func(r rune) bool { return r == ' ' || r == ',' || r == '\t' }
		for _, s := range strings.FieldsFunc(rest, sep) {
			to, err := strconv.Atoi(s)
			if err != nil || to < 0 || to >= n || to == from {
				continue
			}
			a, b := from, to
			if a > b {
				a, b = b, a
			}
			key := [2]int{a, b}
			if seen[key] {
				continue
			}
			seen[key] = true
			g.AddEdge(from, to)
		}
	}
	return g
}
