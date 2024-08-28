package gamelogic

func GetRandomQuestion() string {
	return AllQuestions[0]
}

var AllQuestions = []string{
	"question A",
	"question B"}
