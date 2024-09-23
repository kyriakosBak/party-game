package gamelogic

import (
	"log/slog"
	"math/rand"
	"strings"
)

var shuffledQuestions []string
var playerNamePlaceholder string = "[player's name]"

func GetRandomQuestion(playerName string) string {
	if shuffledQuestions == nil {
		slog.Debug("Shuffling questions")
		shuffledQuestions = shuffle(questions)
	}

	result := shuffledQuestions[len(shuffledQuestions)-1]
	shuffledQuestions = shuffledQuestions[:len(shuffledQuestions)-1]
	result = strings.Replace(result, playerNamePlaceholder, playerName, 1)

	return result
}

// Shuffle function using Fisher-Yates shuffle algorithm
func shuffle(questions []string) []string {
	shuffled := make([]string, len(questions))
	copy(shuffled, questions) // Copy the original list to avoid modifying it

	// Fisher-Yates shuffle
	for i := len(shuffled) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return shuffled
}

var questions = []string{
	"What’s the strangest thing [player's name] has ever eaten?",
	"If [player's name] could only eat one food for the rest of their life, what would it be?",
	"What’s [player's name] biggest guilty pleasure?",
	"What song does [player's name] secretly love but would never admit?",
	"What is [player's name] most embarrassing moment?",
	"If [player's name] could have any superpower, what would it be?",
	"What is [player's name] irrational fear?",
	"What’s [player's name] favorite movie of all time?",
	"What’s one thing [player's name] can’t live without?",
	"If [player's name] could swap lives with any celebrity for a day, who would it be?",
	"What is the worst job [player's name] has ever had?",
	"If [player's name] won the lottery, what’s the first thing they would buy?",
	"What’s [player's name] go-to comfort food?",
	"What is [player's name] worst habit?",
	"What is [player's name] biggest pet peeve?",
	"If [player's name] could time travel, which historical event would they visit?",
	"What’s one thing [player's name] is surprisingly good at?",
	"What was [player's name] favorite childhood toy?",
	"If [player's name] could live anywhere in the world, where would it be?",
	"What is [player's name] hidden talent?",
	"What’s [player's name] dream vacation?",
	"What is [player's name] least favorite household chore?",
	"What’s the craziest thing [player's name] has done for love?",
	"What does [player's name] always have in their fridge?",
	"What’s [player's name] favorite hobby?",
	"If [player's name] could be any fictional character, who would it be?",
	"What’s the most embarrassing thing in [player's name] search history?",
	"What’s the weirdest thing [player's name] has ever Googled?",
	"If [player's name] had to describe themselves in one word, what would it be?",
	"What’s the longest [player's name] has gone without sleep?",
	"If [player's name] were famous, what would they be famous for?",
	"What’s [player's name] favorite way to spend a lazy day?",
	"What is [player's name] dream job?",
	"What’s the most rebellious thing [player's name] did as a teenager?",
	"What is [player's name] biggest regret?",
	"If [player's name] could instantly master any skill, what would it be?",
	"What’s the worst haircut [player's name] has ever had?",
	"What’s the most unusual compliment [player's name] has ever received?",
	"If [player's name] could have dinner with any historical figure, who would it be?",
	"What’s [player's name] favorite childhood memory?",
	"What’s the funniest thing [player's name] has ever witnessed?",
	"What’s the worst gift [player's name] has ever received?",
	"If [player's name] had a theme song, what would it be?",
	"What is [player's name] most-used emoji?",
	"What is [player's name] go-to karaoke song?",
	"If [player's name] could redo one day in their life, which day would it be?",
	"What’s the weirdest job [player's name] has ever had?",
	"What was the first concert [player's name] ever went to?",
	"What’s the most ridiculous thing [player's name] has ever spent money on?",
	"What’s [player's name] worst fashion choice?",
	"What would [player's name] name their autobiography?",
	"What’s the most ridiculous law [player's name] would create if they were a dictator?",
	"What would be [player's name] go-to excuse if they were late to their own wedding?",
	"What’s the weirdest thing [player's name] would do if they were invisible for a day?",
	"If [player's name] were a superhero, what would their catchphrase be?",
	"What unusual career would [player's name] choose if they weren’t doing what they do now?",
	"What would [player's name] do if they woke up one day as the opposite gender?",
	"If [player's name] were an alien, what would their home planet be like?",
	"What is the most absurd rumor [player's name] could start about themselves?",
	"If [player's name] were an animal, what would they be and why?",
	"What outlandish invention would [player's name] create to solve world problems?",
	"What would be [player's name] strategy to survive a zombie apocalypse?",
	"What’s the weirdest thing [player's name] would put on their bucket list?",
	"If [player's name] were a reality TV star, what would their show be about?",
	"If [player's name] could erase one thing from existence, what would it be?",
	"What would [player's name] do if they were stuck in an elevator with their worst enemy?",
	"What ridiculous thing would [player's name] do for a million dollars?",
	"What would be [player's name] most bizarre world record?",
	"If [player's name] were a flavor of ice cream, what would they be?",
	"What would [player's name] do if they had to live in a world without internet?",
	"What is the craziest conspiracy theory [player's name] secretly believes?",
	"What is the silliest superstition [player's name] has?",
	"What embarrassing thing would [player's name] likely be caught doing in public?",
	"If [player's name] could only say one word for the rest of their life, what would it be?",
	"What would [player's name] spend a day doing if there were no consequences?",
	"If [player's name] could instantly switch lives with any fictional villain, who would it be?",
	"What fictional place would [player's name] choose to live in?",
	"If [player's name] had to get a tattoo of something ridiculous, what would it be?",
	"What would [player's name] rename the planet if given the chance?",
	"If [player's name] could replace one body part with a robotic one, which would it be?",
	"What would [player's name] do if they woke up in the wrong decade?",
	"If [player's name] could speak to animals, what would they ask first?",
	"What would be [player's name] ultimate prank?",
	"If [player's name] could create their own holiday, what would it celebrate?",
	"What would [player's name] sell if they opened a shop?",
	"What mythical creature would [player's name] most likely have as a pet?",
	"What is the weirdest dream [player's name] has ever had?",
	"If [player's name] could teleport anywhere, where would they go first?",
	"What bizarre item would [player's name] take with them to a deserted island?",
	"If [player's name] could swap lives with an animal for a day, what animal would they choose?",
	"What outlandish talent show would [player's name] be a contestant on?",
	"If [player's name] could time travel to any era, where would they go and why?",
	"What would be [player's name] strategy to win a reality competition show?",
	"What would be [player's name] first act as ruler of the world?",
	"If [player's name] were to host a cooking show, what strange dish would they cook?",
	"What is [player's name] most absurd party trick?",
	"What cartoon character is most like [player's name]?",
	"If [player's name] could own any fictional weapon, what would it be?",
	"What’s the weirdest nickname [player's name] could give themselves?",
	"What alien planet would [player's name] most likely want to visit?",
	"What strange hobby would [player's name] take up if they had unlimited free time?",
	"If [player's name] were a wizard, what ridiculous spell would they invent?",
	"What would [player's name] name their spaceship?",
	"What’s [player's name] least favorite fictional character?",
	"What’s the most ridiculous goal [player's name] could set for themselves?",
	"If [player's name] could only wear one color for the rest of their life, what would it be?",
	"If [player's name] could swap brains with someone for a day, who would it be?",
	"What strange talent would [player's name] gain from a magic spell gone wrong?",
	"What would [player's name] do if they found a hidden door in their house?",
	"What would [player's name] wish for if they found a genie but could only ask for weird things?",
}
