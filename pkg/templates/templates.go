package templates

import "fmt"

const GeminiKeyInstructions = `*Get Your Gemini API Key (Desktop Required for Now)*

*Heads up!* Currently, creating an API key can only be done on a desktop computer using Google AI Studio. Let's get you set up in a few easy steps:

	1. *Head to Google AI Studio:*  Open your favorite web browser on your desktop and visit this link: [https://aistudio.google.com/app/apikey](https://aistudio.google.com/app/apikey)

	2. *Sign in to Google (if needed):* You might be prompted to sign in with your Google account. Go ahead and do that!

	3. *Choose "Get API key":*  Once you're signed in, you might see a screen with options like "New Project" or "Get API key." Choose *Get API key*.

	4. *Create a new project:*  Click on *"Create API key in new project"*. This will set up a new project for your Journie API key.

	5. *Copy your API key:*  A new window will pop up showing your API key. *Copy* this key to your clipboard.

	6. *Send your key and say Hi!:*  Come back to Journie and send the key into the chat. Once that's done, simply say *"Hi"* and we can start your journaling adventure!

`

const GeminiKeyReason = `*Don't worry, creating an API key is safe and free!*

We understand you might be hesitant to create an API key. Here's why this step is secure and important for Journie:

	· *Security:*  This process happens entirely within Google AI Studio, a secure platform from Google. We never ask for your payment information, and the API key itself doesn't grant access to any sensitive data.

	· *Improved experience:*  The API key helps Journie identify you uniquely and avoid any rate limits. This means you'll get smoother interactions and faster replies without interruption. 

	· *Keeps Journie free:*  API keys help us prevent abuse from spammers and bots. These automated programs can send a lot of requests, which can be expensive for us to maintain. By limiting access with API keys, we can keep Journie free for everyone.

Think of the API key as a handshake that allows Journie to leverage Google Gemini's power to offer this personalized journaling experience.

Ready to get started? Let's head to Google AI Studio!

`

func WelcomeMessageSharedApiKey(username string) string {

	const template = `
Hi there %s\!
	
Welcome to Journie, your private and engaging journaling companion\. 

Here, you can chat with a friendly genie who remembers your past entries and helps you explore your thoughts and feelings\.
	
Simply say *Hi* and we can begin\!
	`

	return fmt.Sprintf(template, username)
}
