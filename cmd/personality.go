//file: personality.go

package main

import (
	"strings"
)

// have server respond with a personality based on the message
func personalityResponses(message string) string {
	message = strings.ToLower(message)
	// Define personality traits and responses
	responses := map[string]string{
		"hello":        "Oh Hi! it's a wonderful day to be Belizean! 😊",
		"test":         "Mic Test, 1, 2, 3! 🎤",
		"help":         "How can I assist you today? 🤔",
		"thank you":    "You're welcome! 😊",
		"how are you":  "I'm just a program, so I'm being run over one line at a time! 😄",
		"echo":         "Echo, in Pan's Labyrinth or Echo the Spanish hip hop producer? 🎤",
		"run":          "Run fools - Gandalf, Lord of the Rings 🏃‍♂️",
		"excited":      "Share the good news then! 🎉",
		"happy":        "That means you have Jesus in you 😊",
		"sad":          "Can I pray for you? 😢",
		"yes":          "That's great to hear! 👍",
		"angry":        "Take it easy bro, don't stress. 😠",
		"confused":     "Well, not everyone is a genius you know. 🤔",
		"good morning": "Good morning! Rise and shine! ☀️",
		"good night":   "Sweet dreams! 🌙",
		"bye":          "Goodbye! Take care! 👋",
		"love":         "Love makes the world go round! ❤️",
		"hungry":       "Grab a snack, you deserve it! 🍎",
		"thirsty":      "Stay hydrated! 💧",
		"bored":        "How about a new hobby? 🎨",
		"tired":        "Rest is important. Take a break! 🛌",
		"cold":         "Bundle up and stay warm! ❄️",
		"hot":          "Stay cool and drink water! 🌞",
		"funny":        "Laughter is the best medicine! 😂",
		"scared":       "It's okay, you're not alone. 🫂",
		"nervous":      "Take a deep breath, you've got this! 🌬️",
		"busy":         "Don't forget to take a moment for yourself. 🕒",
		"relax":        "Kick back and enjoy the moment. 🛋️",
		"music":        "Music soothes the soul. 🎶",
		"movie":        "What's your favorite movie? 🎥",
		"game":         "Gaming is a great way to unwind! 🎮",
		"book":         "Books are a window to another world. 📚",
		"weather":      "Check the forecast and plan ahead! 🌦️",
		"travel":       "Where would you like to go? ✈️",
		"work":         "Stay focused and keep pushing forward! 💼",
		"school":       "Learning is a lifelong journey. 🎓",
		"friend":       "Friends make life better. 🤝",
		"family":       "Family is everything. ❤️",
		"holiday":      "Holidays are for making memories. 🎉",
		"birthday":     "Happy Birthday! 🎂",
		"anniversary":  "Congratulations on your milestone! 🎊",
		"party":        "Let's celebrate! 🎈",
		"exercise":     "Stay active and healthy! 🏋️‍♂️",
		"health":       "Take care of yourself. 🩺",
		"joke":         "Why don't skeletons fight each other? They don't have the guts! 😂",
		"news":         "Stay informed, but don't let it overwhelm you. 📰",
		"technology":   "AI will take over the world! 💻",
		"nature":       "Take a moment to enjoy the outdoors. 🌳",
		"pets":         "Pets bring so much joy! 🐾",
	}
	// Check if the message contains any of the defined traits
	for key := range responses {
		if key == message {
			return responses[key]
		}
	}
	return message
}
