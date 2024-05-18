import Genie from "./components/genie";
import QR from "./components/qr";

export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-between p-6 lg:p-16">
      <div className="flex flex-col items-center justify-between mb-8 lg:mb-12 before:fixed before:h-[300px] before:w-full before:-translate-x-1/2 before:rounded-full before:bg-gradient-radial before:from-white before:to-transparent before:blur-2xl before:content-[''] after:fixed after:-z-20 after:h-[180px] after:w-full after:translate-x-1/3 after:bg-gradient-conic after:from-sky-200 after:via-blue-200 after:blur-2xl after:content-[''] before:dark:bg-gradient-to-br before:dark:from-transparent before:dark:to-blue-700 before:dark:opacity-10 after:dark:from-sky-900 after:dark:via-[#0141ff] after:dark:opacity-40 sm:before:w-[480px] sm:after:w-[240px] before:lg:h-[360px]">
        <div className="sm:hidden">
          <Genie size={"sm"} />
        </div>

        <div className="max-sm:hidden">
          <Genie size={"lg"} />
        </div>

        <div className="text-[44px] lg:text-[168px] -mt-5 lg:-mt-20 font-bold  animate-text bg-gradient-to-r from-teal-500 via-purple-500 to-orange-500 bg-clip-text text-transparent">
          Journie
          <span className="text-base font-medium align-text-top	">Alpha</span>
        </div>

        <div className="font-normal text-[20px] lg:text-[27px] text-center animate-text bg-gradient-to-r from-orange-500 via-purple-500 to-teal-500 bg-clip-text text-transparent">
          Your AI-powered Journaling Genie on Telegram
        </div>
      </div>

      <div className="mb-8 lg:mb-12 grid text-center w-full lg:max-w-6xl lg:grid-cols-4 lg:text-left">
        <div className="group rounded-lg border border-transparent px-5 py-4 transition-colors hover:border-gray-300 hover:bg-gray-100 hover:dark:border-neutral-700 hover:dark:bg-neutral-800/30">
          <h2 className="mb-3 text-xl lg:text-3xl font-semibold">
            Personal{" "}
            <span className="inline-block transition-transform group-hover:translate-x-1 motion-reduce:transform-none"></span>
          </h2>
          <p className="m-0 opacity-50">
            Journie remembers your conversations, turning them into insightful
            journaling journeys with personalized prompts.
          </p>
        </div>

        <div className="group rounded-lg border border-transparent px-5 py-4 transition-colors hover:border-gray-300 hover:bg-gray-100 hover:dark:border-neutral-700 hover:dark:bg-neutral-800/30">
          <h2 className="mb-3 text-xl lg:text-3xl font-semibold">
            Hassle-free{" "}
            <span className="inline-block transition-transform group-hover:translate-x-1 motion-reduce:transform-none"></span>
          </h2>
          <p className="m-0 opacity-50">
            Journie removes the friction of starting a new habit – simply open
            Telegram and start writing. No sign up required.
          </p>
        </div>

        <div className="group rounded-lg border border-transparent px-5 py-4 transition-colors hover:border-gray-300 hover:bg-gray-100 hover:dark:border-neutral-700 hover:dark:bg-neutral-800/30">
          <h2 className="mb-3 text-xl lg:text-3xl font-semibold">
            Private{" "}
            <span className="inline-block transition-transform group-hover:translate-x-1 motion-reduce:transform-none"></span>
          </h2>
          <p className="m-0 opacity-50">
            Journie prioritizes your privacy. Sensitive data like email and
            phone numbers are not stored.{" "}
          </p>
        </div>

        <div className="group rounded-lg border border-transparent px-5 py-4 transition-colors hover:border-gray-300 hover:bg-gray-100 hover:dark:border-neutral-700 hover:dark:bg-neutral-800/30">
          <h2 className="mb-3 text-xl lg:text-3xl font-semibold">
            Gratis{" "}
            <span className="inline-block transition-transform group-hover:translate-x-1 motion-reduce:transform-none"></span>
          </h2>
          <p className="m-0 text-balance opacity-50">
            Start your journaling journey for free! Journie is currently
            available at no cost.{" "}
          </p>
        </div>
      </div>

      <div className="mb-12 lg:mb-16 flex flex-col items-center justify-between gap-4">
        <div className="py-4 lg:py-8 text-center font-medium text-3xl lg:text-[48px] animate-text bg-gradient-to-r from-teal-500 via-purple-500 to-orange-500 bg-clip-text text-transparent">
          Start your Journie
        </div>

        <p className="font-normal text-lg">
          Scan QR code below or tap to open in Telegram
        </p>

        <div className="p-8">
          <div className="sm:hidden">
            <QR size={"sm"} />
          </div>

          <div className="max-sm:hidden">
            <QR size={"lg"} />
          </div>
        </div>

        <div className="p-8 text-2xl lg:text-3xl font-semibold">
          Happy Journaling! 🎉
        </div>
      </div>

      <div className="mt-8 flex flex-col items-center justify-between mb-32 gap-4">
        <div className="py-4 lg:py-8 text-center font-medium text-3xl lg:text-[48px] animate-text bg-gradient-to-r  from-orange-500 via-purple-500 to-teal-500 bg-clip-text text-transparent">
          Still Curious?
        </div>
        <p className="font-normal text-lg">
          Here are all the answers (probably)
        </p>

        <div className="mb-28 grid lg:mb-0 lg:w-full lg:max-w-5xl lg:grid-cols-1 lg:text-left">
          <div className="min-h-32 group rounded-lg border border-transparent px-5 py-4 transition-colors hover:border-gray-300 hover:bg-gray-100 hover:dark:border-neutral-700 hover:dark:bg-neutral-800/30">
            <h2 className="mb-3 text-xl lg:text-2xl font-semibold">
              Who{" "}
              <span className=" inline-block opacity-35">
                is behind Journie
              </span>
            </h2>
            <p className="m-0">
              <span className="m-0 opacity-50">
                Full-stack Dev from Singapore, keyboard mashing for 8 years and
                counting. Connect on LinkedIn{" "}
              </span>
              <a
                href="https://www.linkedin.com/in/ngping/"
                className="text-blue-500"
                target="_blank"
              >
                Here!
              </a>
            </p>
          </div>

          <div className="min-h-32 group rounded-lg border border-transparent px-5 py-4 transition-colors hover:border-gray-300 hover:bg-gray-100 hover:dark:border-neutral-700 hover:dark:bg-neutral-800/30">
            <h2 className="mb-3 text-2xl font-semibold">
              What <span className=" inline-block opacity-35">is Journie</span>
            </h2>
            <p className="m-0 opacity-50">
              Powered by Google’s Gemini AI, Journie integrates with Telegram
              Bot API and glued together with Go and Firebase. Website is made
              with NextJS.
            </p>
          </div>
          <div className="min-h-32 group rounded-lg border border-transparent px-5 py-4 transition-colors hover:border-gray-300 hover:bg-gray-100 hover:dark:border-neutral-700 hover:dark:bg-neutral-800/30">
            <h2 className="mb-3 text-2xl font-semibold">
              When{" "}
              <span className=" inline-block opacity-35">
                did Journie start
              </span>
            </h2>
            <p className="m-0 opacity-50">
              Journie&apos;s journey began Spring 2024.
            </p>
          </div>
          <div className="min-h-32 group rounded-lg border border-transparent px-5 py-4 transition-colors hover:border-gray-300 hover:bg-gray-100 hover:dark:border-neutral-700 hover:dark:bg-neutral-800/30">
            <h2 className="mb-3 text-2xl font-semibold">
              Where{" "}
              <span className=" inline-block opacity-35">
                is Journie headed
              </span>
            </h2>
            <p className="m-0 opacity-50">
              Scale to support more users while fine-tuning the model to address
              themes like mental health, productivity, sleep, and even exploring
              paid tiers to further enhance the experience.
            </p>
          </div>
          <div className="min-h-32 group rounded-lg border border-transparent px-5 py-4 transition-colors hover:border-gray-300 hover:bg-gray-100 hover:dark:border-neutral-700 hover:dark:bg-neutral-800/30">
            <h2 className="mb-3 text-2xl font-semibold">
              Why <span className=" inline-block opacity-35">Journie</span>
            </h2>
            <p className="m-0 opacity-50">
              Started out as a simple Telegram bot bedtime reminders, it evolved
              into a Journaling tool to help the creator (and you!) get into the
              mindful habit of journaling with minimal hassle and maximum
              engagement.
            </p>
          </div>
        </div>
      </div>
    </main>
  );
}
