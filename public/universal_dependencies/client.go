package universaldependencies

// UniversalDependecies
type UniversalDependencies struct {
	internalUniversalDependecies map[string]string
}

func NewUniversalDependencies() *UniversalDependencies {
	internalUniversalDependecies := make(map[string]string)
	internalUniversalDependecies["amod"] = "определение"
	internalUniversalDependecies["nsubj"] = "номинальный_субъект" // подлежащее
	internalUniversalDependecies["obj"] = "объект"                // дополнение
	internalUniversalDependecies["punct"] = "знак_пунктуации"
	internalUniversalDependecies["root"] = "основа"          // может быть глагол
	internalUniversalDependecies["obl"] = "объект_локальный" //дополнение
	internalUniversalDependecies["case"] = "указатель"       // предлог
	internalUniversalDependecies["fixed"] = "фиксированный"
	internalUniversalDependecies["det"] = "определитель" // детерминатив
	internalUniversalDependecies["parataxis"] = "сочиненный"
	internalUniversalDependecies["iobj"] = "косвенное_дополнение"
	internalUniversalDependecies["nmod"] = "номинальный_модификатор"
	internalUniversalDependecies["advmod"] = "наречный_модификатор"
	internalUniversalDependecies["conj"] = "соединение"
	internalUniversalDependecies["nsubj:pass"] = "пассивный_номинальный_субъект"
	internalUniversalDependecies["aux:pass"] = "пассивный_вспомогательный"
	internalUniversalDependecies["aux"] = "вспомогательный"
	internalUniversalDependecies["xcomp"] = "открытое_клаузальное_дополнение"
	internalUniversalDependecies["cop"] = "связка"
	internalUniversalDependecies["advcl"] = "модификатор_придаточного_предложения"
	internalUniversalDependecies["mark"] = "маркер"
	internalUniversalDependecies["expl"] = "указательная_частица" // _как_«лишнее»_подлежащее
	internalUniversalDependecies["cc"] = "координирующее_соединение"
	internalUniversalDependecies["acl:relcl"] = "модификатор_относительного_предложения"
	internalUniversalDependecies["csubj"] = "клаузальный_субъект"
	internalUniversalDependecies["acl"] = "клаузальный_модификатор_существительного" // _придаточного_предложения
	internalUniversalDependecies["nummod:gov"] = "числовой_модификатор_регулирующий_падеж_существительного"
	internalUniversalDependecies["ccomp"] = "клаузальное_дополнение"
	internalUniversalDependecies["appos"] = "аппозиционный_модификатор_существительного"
	return &UniversalDependencies{
		internalUniversalDependecies: internalUniversalDependecies,
	}
}

func (ud *UniversalDependencies) Tag2StrByUniversalDependecies(tag string) string {
	return ud.internalUniversalDependecies[tag]
}
