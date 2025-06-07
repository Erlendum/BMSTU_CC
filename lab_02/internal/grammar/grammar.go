package grammar

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Grammar struct {
	NonTerminals    []string
	Terminals       []string
	ProductionRules []ProductionRule
	StartSymbol     string
	productionMap   map[string][]string
}

type ProductionRule struct {
	Left  string
	Right string
}

func NewGrammarFromString(input string) (*Grammar, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) < 7 {
		return nil, fmt.Errorf("недостаточно строк во входных данных (ожидается 7, получено %d)", len(lines))
	}

	var cleanLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			cleanLines = append(cleanLines, trimmed)
		}
	}

	g := &Grammar{productionMap: make(map[string][]string)}

	nonTerminalCount, err := strconv.Atoi(cleanLines[0])
	if err != nil {
		return nil, fmt.Errorf("неверный формат числа нетерминалов: %v", err)
	}

	g.NonTerminals = strings.Fields(cleanLines[1])
	if len(g.NonTerminals) != nonTerminalCount {
		return nil, fmt.Errorf("несоответствие количества нетерминалов: ожидается %d, получено %d",
			nonTerminalCount, len(g.NonTerminals))
	}

	terminalCount, err := strconv.Atoi(cleanLines[2])
	if err != nil {
		return nil, fmt.Errorf("неверный формат числа терминалов: %v", err)
	}

	g.Terminals = strings.Fields(cleanLines[3])
	if len(g.Terminals) != terminalCount {
		return nil, fmt.Errorf("несоответствие количества терминалов: ожидается %d, получено %d",
			terminalCount, len(g.Terminals))
	}

	ruleCount, err := strconv.Atoi(cleanLines[4])
	if err != nil {
		return nil, fmt.Errorf("неверный формат числа правил вывода: %v", err)
	}

	ruleStart := 5
	ruleEnd := ruleStart + ruleCount
	if ruleEnd > len(cleanLines)-1 {
		return nil, fmt.Errorf("недостаточно правил вывода: ожидается %d, доступно %d",
			ruleCount, len(cleanLines)-ruleStart-1)
	}

	for i := ruleStart; i < ruleEnd; i++ {
		parts := strings.Split(cleanLines[i], "->")
		if len(parts) != 2 {
			return nil, fmt.Errorf("неправильный формат правила вывода: %s", cleanLines[i])
		}
		rule := ProductionRule{
			Left:  strings.TrimSpace(parts[0]),
			Right: strings.TrimSpace(parts[1]),
		}
		g.ProductionRules = append(g.ProductionRules, rule)
		g.productionMap[rule.Left] = append(g.productionMap[rule.Left], rule.Right)
	}

	g.StartSymbol = strings.TrimSpace(cleanLines[ruleEnd])
	found := false
	for _, nt := range g.NonTerminals {
		if nt == g.StartSymbol {
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("начальный символ %s не найден в нетерминалах", g.StartSymbol)
	}

	return g, nil
}

func (g *Grammar) ToString() string {
	var builder strings.Builder

	sortedNonTerminals := append([]string{}, g.NonTerminals...)
	sort.Strings(sortedNonTerminals)
	builder.WriteString(fmt.Sprintf("%d\n", len(sortedNonTerminals)))
	builder.WriteString(strings.Join(sortedNonTerminals, " ") + "\n")

	sortedTerminals := append([]string{}, g.Terminals...)
	sort.Strings(sortedTerminals)
	builder.WriteString(fmt.Sprintf("%d\n", len(sortedTerminals)))
	builder.WriteString(strings.Join(sortedTerminals, " ") + "\n")

	rulesCopy := append([]ProductionRule{}, g.ProductionRules...)
	sort.Slice(rulesCopy, func(i, j int) bool {
		if rulesCopy[i].Left == rulesCopy[j].Left {
			return rulesCopy[i].Right < rulesCopy[j].Right
		}
		return rulesCopy[i].Left < rulesCopy[j].Left
	})

	builder.WriteString(fmt.Sprintf("%d\n", len(rulesCopy)))
	for _, rule := range rulesCopy {
		builder.WriteString(fmt.Sprintf("%s -> %s\n", rule.Left, rule.Right))
	}

	builder.WriteString(g.StartSymbol + "\n")

	return builder.String()
}

func (g *Grammar) GetProductionsFor(nonTerminal string) []string {
	return g.productionMap[nonTerminal]
}

func (g *Grammar) IsNonTerminal(symbol string) bool {
	for _, nt := range g.NonTerminals {
		if nt == symbol {
			return true
		}
	}
	return false
}

func (g *Grammar) IsTerminal(symbol string) bool {
	for _, t := range g.Terminals {
		if t == symbol {
			return true
		}
	}
	return false
}

func (g *Grammar) Print() {
	fmt.Println("Нетерминалы:", strings.Join(g.NonTerminals, ", "))
	fmt.Println("Терминалы:", strings.Join(g.Terminals, ", "))
	fmt.Println("Начальный символ:", g.StartSymbol)
	fmt.Println("Правила вывода:")
	for _, rule := range g.ProductionRules {
		fmt.Printf("  %s -> %s\n", rule.Left, rule.Right)
	}
}

func (g *Grammar) EliminateLeftRecursion() error {
	orderedNT := g.NonTerminals

	for i, Ai := range orderedNT {
		for j := 0; j < i; j++ {
			var productionsToAdd []ProductionRule
			var productionsToRemove []ProductionRule
			Ai := orderedNT[i]
			Aj := orderedNT[j]
			for _, rule := range g.ProductionRules {
				if strings.HasPrefix(rule.Left, Ai) && strings.HasPrefix(rule.Right, Aj) {
					gamma := rule.Right[len(Aj):]
					for _, ajRule := range g.GetProductionsFor(Aj) {
						newRight := ajRule + gamma
						productionsToAdd = append(productionsToAdd, ProductionRule{
							Left:  Ai,
							Right: newRight,
						})
					}
					productionsToRemove = append(productionsToRemove, rule)
					g.removeProductions(productionsToRemove)
					g.ProductionRules = append(g.ProductionRules, productionsToAdd...)
					g.rebuildProductionMap()
					productionsToAdd = []ProductionRule{}
					productionsToRemove = []ProductionRule{}
				}
			}
		}
		if err := g.eliminateImmediateLeftRecursion(Ai); err != nil {
			return fmt.Errorf("ошибка устранения непосредственной левой рекурсии для %s: %v", Ai, err)
		}
	}

	return nil
}

func (g *Grammar) eliminateImmediateLeftRecursion(A string) error {
	var recursiveRules []ProductionRule
	var nonRecursiveRules []ProductionRule

	for _, rule := range g.GetProductionsFor(A) {
		if strings.HasPrefix(rule, A) {
			recursiveRules = append(recursiveRules, ProductionRule{
				Left:  A,
				Right: rule,
			})
		} else {
			nonRecursiveRules = append(nonRecursiveRules, ProductionRule{
				Left:  A,
				Right: rule,
			})
		}
	}

	if len(recursiveRules) == 0 {
		return nil
	}

	newNonTerminal := A + "'"
	g.NonTerminals = append(g.NonTerminals, newNonTerminal)

	var newARules []ProductionRule
	for _, rule := range nonRecursiveRules {
		if rule.Right == "ε" {
			newARules = append(newARules, ProductionRule{
				Left:  A,
				Right: newNonTerminal,
			})
		} else {
			newARules = append(newARules, ProductionRule{
				Left:  A,
				Right: rule.Right + newNonTerminal,
			})
		}
	}

	var newAPrimeRules []ProductionRule
	for _, rule := range recursiveRules {
		alpha := rule.Right[len(A):]
		newAPrimeRules = append(newAPrimeRules, ProductionRule{
			Left:  newNonTerminal,
			Right: alpha + newNonTerminal,
		})
	}
	newAPrimeRules = append(newAPrimeRules, ProductionRule{
		Left:  newNonTerminal,
		Right: "ε",
	})

	g.removeProductionsFor(A)
	g.ProductionRules = append(g.ProductionRules, newARules...)
	g.ProductionRules = append(g.ProductionRules, newAPrimeRules...)
	g.rebuildProductionMap()

	return nil
}

func (g *Grammar) removeProductions(rules []ProductionRule) {
	for _, rule := range rules {
		for i, r := range g.ProductionRules {
			if r.Left == rule.Left && r.Right == rule.Right {
				g.ProductionRules = append(g.ProductionRules[:i], g.ProductionRules[i+1:]...)
				break
			}
		}
	}
}

func (g *Grammar) removeProductionsFor(nt string) {
	var newRules []ProductionRule
	for _, rule := range g.ProductionRules {
		if rule.Left != nt {
			newRules = append(newRules, rule)
		}
	}
	g.ProductionRules = newRules
}

func (g *Grammar) rebuildProductionMap() {
	g.productionMap = make(map[string][]string)
	for _, rule := range g.ProductionRules {
		g.productionMap[rule.Left] = append(g.productionMap[rule.Left], rule.Right)
	}
}

func (g *Grammar) EliminateChainRules() error {
	Nsets := make(map[string]map[string]bool)

	for _, A := range g.NonTerminals {
		Nprev := make(map[string]bool)
		Nprev[A] = true
		changed := true

		for changed {
			changed = false
			Ncurrent := make(map[string]bool)
			for B := range Nprev {
				Ncurrent[B] = true
				for _, rule := range g.ProductionRules {
					if rule.Left == B && g.IsNonTerminal(rule.Right) && len(rule.Right) == 1 {
						C := rule.Right
						if !Nprev[C] {
							Ncurrent[C] = true
							changed = true
						}
					}
				}
			}
			Nprev = Ncurrent
		}

		Nsets[A] = Nprev
	}

	newRulesMap := make(map[ProductionRule]bool)

	for _, rule := range g.ProductionRules {
		if rule.Right == "ε" || (g.IsNonTerminal(rule.Right) && len(rule.Right) == 1) {
			continue
		}

		for A, N_A := range Nsets {
			if N_A[rule.Left] {
				newRule := ProductionRule{Left: A, Right: rule.Right}
				newRulesMap[newRule] = true
			}
		}
	}

	var newRules []ProductionRule
	for rule := range newRulesMap {
		newRules = append(newRules, rule)
	}

	g.ProductionRules = newRules
	g.rebuildProductionMap()

	return nil
}
