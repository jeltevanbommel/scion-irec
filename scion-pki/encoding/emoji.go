// Copyright 2022 Anapaya Systems
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package encoding

import "strings"

func ToEmoji(input []byte) string {
	var b strings.Builder
	for _, r := range input {
		b.WriteString(emojiCodeMap[r])
	}
	return b.String()
}

// Encoding mapping for emojis. This mapping is taken from smallstep to provide
// compatibility
// (https://github.com/smallstep/cli/blob/1e0b1667db00f1b0756e0d833fca496a323338b7/crypto/fingerprint/emoji.go).
var emojiCodeMap = []string{
	"\U0001f44d",                   // 👍 :+1:
	"\U0001f3b1",                   // 🎱 :8ball:
	"\u2708\ufe0f",                 // ✈️ :airplane:
	"\U0001f47d",                   // 👽 :alien:
	"\u2693",                       // ⚓ :anchor:
	"\U0001f47c",                   // 👼 :angel:
	"\U0001f620",                   // 😠 :angry:
	"\U0001f41c",                   // 🐜 :ant:
	"\U0001f34e",                   // 🍎 :apple:
	"\U0001f3a8",                   // 🎨 :art:
	"\U0001f476",                   // 👶 :baby:
	"\U0001f37c",                   // 🍼 :baby_bottle:
	"\U0001f519",                   // 🔙 :back:
	"\U0001f38d",                   // 🎍 :bamboo:
	"\U0001f34c",                   // 🍌 :banana:
	"\U0001f488",                   // 💈 :barber:
	"\U0001f6c1",                   // 🛁 :bathtub:
	"\U0001f37a",                   // 🍺 :beer:
	"\U0001f514",                   // 🔔 :bell:
	"\U0001f6b4\u200d\u2642\ufe0f", // 🚴‍♂️ :bicyclist:
	"\U0001f426",                   // 🐦 :bird:
	"\U0001f382",                   // 🎂 :birthday:
	"\U0001f33c",                   // 🌼 :blossom:
	"\U0001f699",                   // 🚙 :blue_car:
	"\U0001f417",                   // 🐗 :boar:
	"\U0001f4a3",                   // 💣 :bomb:
	"\U0001f4a5",                   // 💥 :boom:
	"\U0001f647\u200d\u2642\ufe0f", // 🙇‍♂️ :bow:
	"\U0001f466",                   // 👦 :boy:
	"\U0001f494",                   // 💔 :broken_heart:
	"\U0001f4a1",                   // 💡 :bulb:
	"\U0001f68c",                   // 🚌 :bus:
	"\U0001f335",                   // 🌵 :cactus:
	"\U0001f4c6",                   // 📆 :calendar:
	"\U0001f4f7",                   // 📷 :camera:
	"\U0001f36c",                   // 🍬 :candy:
	"\U0001f431",                   // 🐱 :cat:
	"\U0001f352",                   // 🍒 :cherries:
	"\U0001f6b8",                   // 🚸 :children_crossing:
	"\U0001f36b",                   // 🍫 :chocolate_bar:
	"\U0001f44f",                   // 👏 :clap:
	"\u2601\ufe0f",                 // ☁️ :cloud:
	"\u2663\ufe0f",                 // ♣️ :clubs:
	"\U0001f1e8\U0001f1f3",         // 🇨🇳 :cn:
	"\u2615",                       // ☕ :coffee:
	"\U0001f6a7",                   // 🚧 :construction:
	"\U0001f36a",                   // 🍪 :cookie:
	"\u00a9\ufe0f",                 // ©️ :copyright:
	"\U0001f33d",                   // 🌽 :corn:
	"\U0001f42e",                   // 🐮 :cow:
	"\U0001f319",                   // 🌙 :crescent_moon:
	"\U0001f451",                   // 👑 :crown:
	"\U0001f622",                   // 😢 :cry:
	"\U0001f52e",                   // 🔮 :crystal_ball:
	"\u27b0",                       // ➰ :curly_loop:
	"\U0001f46f\u200d\u2640\ufe0f", // 👯‍♀️ :dancers:
	"\U0001f4a8",                   // 💨 :dash:
	"\U0001f1e9\U0001f1ea",         // 🇩🇪 :de:
	"\u2666\ufe0f",                 // ♦️ :diamonds:
	"\U0001f436",                   // 🐶 :dog:
	"\U0001f369",                   // 🍩 :doughnut:
	"\U0001f409",                   // 🐉 :dragon:
	"\U0001f4c0",                   // 📀 :dvd:
	"\U0001f442",                   // 👂 :ear:
	"\U0001f346",                   // 🍆 :eggplant:
	"\U0001f418",                   // 🐘 :elephant:
	"\U0001f51a",                   // 🔚 :end:
	"\u2709",                       // ✉ :envelope:
	"\U0001f1ea\U0001f1f8",         // 🇪🇸 :es:
	"\U0001f440",                   // 👀 :eyes:
	"\U0001f44a",                   // 👊 :facepunch:
	"\U0001f468\u200d\U0001f469\u200d\U0001f466", // 👨‍👩‍👦 :family:
	"\U0001f3a1",                         // 🎡 :ferris_wheel:
	"\U0001f630",                         // 😰 :cold_sweat:
	"\U0001f525",                         // 🔥 :fire:
	"\U0001f386",                         // 🎆 :fireworks:
	"\U0001f4be",                         // 💾 :floppy_disk:
	"\U0001f3c8",                         // 🏈 :football:
	"\U0001f374",                         // 🍴 :fork_and_knife:
	"\U0001f340",                         // 🍀 :four_leaf_clover:
	"\U0001f1eb\U0001f1f7",               // 🇫🇷 :fr:
	"\U0001f35f",                         // 🍟 :fries:
	"\U0001f95c",                         // 🥜 :peanuts:
	"\U0001f595",                         // 🖕 :fu:
	"\U0001f315",                         // 🌕 :full_moon:
	"\U0001f3b2",                         // 🎲 :game_die:
	"\U0001f1ea\U0001f1fa",               // 🇪🇺 :eu:
	"\U0001f48e",                         // 💎 :gem:
	"\U0001f467",                         // 👧 :girl:
	"\U0001f410",                         // 🐐 :goat:
	"\U0001f62c",                         // 😬 :grimacing:
	"\U0001f601",                         // 😁 :grin:
	"\U0001f482\u200d\u2642\ufe0f",       // 💂‍♂️ :guardsman:
	"\U0001f3b8",                         // 🎸 :guitar:
	"\U0001f52b",                         // 🔫 :gun:
	"\U0001f354",                         // 🍔 :hamburger:
	"\U0001f528",                         // 🔨 :hammer:
	"\U0001f439",                         // 🐹 :hamster:
	"\U0001f649",                         // 🙉 :hear_no_evil:
	"\u2764\ufe0f",                       // ❤️ :heart:
	"\U0001f63b",                         // 😻 :heart_eyes_cat:
	"\u2763\ufe0f",                       // ❣️ :heavy_heart_exclamation:
	"\u2714\ufe0f",                       // ✔️ :heavy_check_mark:
	"\U0001f5ff",                         // 🗿 :moyai:
	"\U0001f3ee",                         // 🏮 :izakaya_lantern:
	"\U0001f681",                         // 🚁 :helicopter:
	"\U0001f52a",                         // 🔪 :hocho:
	"\U0001f41d",                         // 🐝 :honeybee:
	"\U0001f434",                         // 🐴 :horse:
	"\U0001f3c7",                         // 🏇 :horse_racing:
	"\u231b",                             // ⌛ :hourglass:
	"\U0001f3e0",                         // 🏠 :house:
	"\U0001f575\ufe0f\u200d\u2640\ufe0f", // 🕵️‍♀️ :female_detective:
	"\U0001f366",                         // 🍦 :icecream:
	"\U0001f47f",                         // 👿 :imp:
	"\U0001f1ee\U0001f1f9",               // 🇮🇹 :it:
	"\U0001f383",                         // 🎃 :jack_o_lantern:
	"\U0001f47a",                         // 👺 :japanese_goblin:
	"\U0001f1ef\U0001f1f5",               // 🇯🇵 :jp:
	"\U0001f511",                         // 🔑 :key:
	"\U0001f48b",                         // 💋 :kiss:
	"\U0001f63d",                         // 😽 :kissing_cat:
	"\U0001f428",                         // 🐨 :koala:
	"\U0001f1f0\U0001f1f7",               // 🇰🇷 :kr:
	"\U0001f34b",                         // 🍋 :lemon:
	"\U0001f484",                         // 💄 :lipstick:
	"\U0001f512",                         // 🔒 :lock:
	"\U0001f36d",                         // 🍭 :lollipop:
	"\U0001f468",                         // 👨 :man:
	"\U0001f341",                         // 🍁 :maple_leaf:
	"\U0001f637",                         // 😷 :mask:
	"\U0001f918",                         // 🤘 :metal:
	"\U0001f52c",                         // 🔬 :microscope:
	"\U0001f4b0",                         // 💰 :moneybag:
	"\U0001f412",                         // 🐒 :monkey:
	"\U0001f5fb",                         // 🗻 :mount_fuji:
	"\U0001f4aa",                         // 💪 :muscle:
	"\U0001f344",                         // 🍄 :mushroom:
	"\U0001f3b9",                         // 🎹 :musical_keyboard:
	"\U0001f3bc",                         // 🎼 :musical_score:
	"\U0001f485",                         // 💅 :nail_care:
	"\U0001f311",                         // 🌑 :new_moon:
	"\u26d4",                             // ⛔ :no_entry:
	"\U0001f443",                         // 👃 :nose:
	"\U0001f39b\ufe0f",                   // 🎛️ :control_knobs:
	"\U0001f529",                         // 🔩 :nut_and_bolt:
	"\u2b55",                             // ⭕ :o:
	"\U0001f30a",                         // 🌊 :ocean:
	"\U0001f44c",                         // 👌 :ok_hand:
	"\U0001f51b",                         // 🔛 :on:
	"\U0001f4e6",                         // 📦 :package:
	"\U0001f334",                         // 🌴 :palm_tree:
	"\U0001f43c",                         // 🐼 :panda_face:
	"\U0001f4ce",                         // 📎 :paperclip:
	"\u26c5",                             // ⛅ :partly_sunny:
	"\U0001f6c2",                         // 🛂 :passport_control:
	"\U0001f43e",                         // 🐾 :paw_prints:
	"\U0001f351",                         // 🍑 :peach:
	"\U0001f427",                         // 🐧 :penguin:
	"\u260e\ufe0f",                       // ☎️ :phone:
	"\U0001f437",                         // 🐷 :pig:
	"\U0001f48a",                         // 💊 :pill:
	"\U0001f34d",                         // 🍍 :pineapple:
	"\U0001f355",                         // 🍕 :pizza:
	"\U0001f448",                         // 👈 :point_left:
	"\U0001f449",                         // 👉 :point_right:
	"\U0001f4a9",                         // 💩 :poop:
	"\U0001f357",                         // 🍗 :poultry_leg:
	"\U0001f64f",                         // 🙏 :pray:
	"\U0001f478",                         // 👸 :princess:
	"\U0001f45b",                         // 👛 :purse:
	"\U0001f4cc",                         // 📌 :pushpin:
	"\U0001f430",                         // 🐰 :rabbit:
	"\U0001f308",                         // 🌈 :rainbow:
	"\u270b",                             // ✋ :raised_hand:
	"\u267b\ufe0f",                       // ♻️ :recycle:
	"\U0001f697",                         // 🚗 :red_car:
	"\u00ae\ufe0f",                       // ®️ :registered:
	"\U0001f380",                         // 🎀 :ribbon:
	"\U0001f35a",                         // 🍚 :rice:
	"\U0001f680",                         // 🚀 :rocket:
	"\U0001f3a2",                         // 🎢 :roller_coaster:
	"\U0001f413",                         // 🐓 :rooster:
	"\U0001f1f7\U0001f1fa",               // 🇷🇺 :ru:
	"\u26f5",                             // ⛵ :sailboat:
	"\U0001f385",                         // 🎅 :santa:
	"\U0001f6f0\ufe0f",                   // 🛰️ :satellite:
	"\U0001f606",                         // 😆 :satisfied:
	"\U0001f3b7",                         // 🎷 :saxophone:
	"\u2702\ufe0f",                       // ✂️ :scissors:
	"\U0001f648",                         // 🙈 :see_no_evil:
	"\U0001f411",                         // 🐑 :sheep:
	"\U0001f41a",                         // 🐚 :shell:
	"\U0001f45e",                         // 👞 :shoe:
	"\U0001f3bf",                         // 🎿 :ski:
	"\U0001f480",                         // 💀 :skull:
	"\U0001f62a",                         // 😪 :sleepy:
	"\U0001f604",                         // 😄 :smile:
	"\U0001f63a",                         // 😺 :smiley_cat:
	"\U0001f60f",                         // 😏 :smirk:
	"\U0001f6ac",                         // 🚬 :smoking:
	"\U0001f40c",                         // 🐌 :snail:
	"\U0001f40d",                         // 🐍 :snake:
	"\u2744\ufe0f",                       // ❄️ :snowflake:
	"\u26bd",                             // ⚽ :soccer:
	"\U0001f51c",                         // 🔜 :soon:
	"\U0001f47e",                         // 👾 :space_invader:
	"\u2660\ufe0f",                       // ♠️ :spades:
	"\U0001f64a",                         // 🙊 :speak_no_evil:
	"\u2b50",                             // ⭐ :star:
	"\u26f2",                             // ⛲ :fountain:
	"\U0001f5fd",                         // 🗽 :statue_of_liberty:
	"\U0001f682",                         // 🚂 :steam_locomotive:
	"\U0001f33b",                         // 🌻 :sunflower:
	"\U0001f60e",                         // 😎 :sunglasses:
	"\u2600\ufe0f",                       // ☀️ :sunny:
	"\U0001f305",                         // 🌅 :sunrise:
	"\U0001f3c4\u200d\u2642\ufe0f",       // 🏄‍♂️ :surfer:
	"\U0001f3ca\u200d\u2642\ufe0f",       // 🏊‍♂️ :swimmer:
	"\U0001f489",                         // 💉 :syringe:
	"\U0001f389",                         // 🎉 :tada:
	"\U0001f34a",                         // 🍊 :tangerine:
	"\U0001f695",                         // 🚕 :taxi:
	"\U0001f3be",                         // 🎾 :tennis:
	"\u26fa",                             // ⛺ :tent:
	"\U0001f4ad",                         // 💭 :thought_balloon:
	"\u2122\ufe0f",                       // ™️ :tm:
	"\U0001f6bd",                         // 🚽 :toilet:
	"\U0001f445",                         // 👅 :tongue:
	"\U0001f3a9",                         // 🎩 :tophat:
	"\U0001f69c",                         // 🚜 :tractor:
	"\U0001f68e",                         // 🚎 :trolleybus:
	"\U0001f922",                         // 🤢 :nauseated_face:
	"\U0001f3c6",                         // 🏆 :trophy:
	"\U0001f3ba",                         // 🎺 :trumpet:
	"\U0001f422",                         // 🐢 :turtle:
	"\U0001f3a0",                         // 🎠 :carousel_horse:
	"\U0001f46d",                         // 👭 :two_women_holding_hands:
	"\U0001f1ec\U0001f1e7",               // 🇬🇧 :uk:
	"\u2602\ufe0f",                       // ☂️ :umbrella:
	"\U0001f513",                         // 🔓 :unlock:
	"\U0001f1fa\U0001f1f8",               // 🇺🇸 :us:
	"\u270c\ufe0f",                       // ✌️ :v:
	"\U0001f4fc",                         // 📼 :vhs:
	"\U0001f3bb",                         // 🎻 :violin:
	"\u26a0\ufe0f",                       // ⚠️ :warning:
	"\U0001f349",                         // 🍉 :watermelon:
	"\U0001f44b",                         // 👋 :wave:
	"\u3030\ufe0f",                       // 〰️ :wavy_dash:
	"\U0001f6be",                         // 🚾 :wc:
	"\u267f",                             // ♿ :wheelchair:
	"\U0001f469",                         // 👩 :woman:
	"\u274c",                             // ❌ :x:
	"\U0001f60b",                         // 😋 :yum:
	"\u26a1",                             // ⚡ :zap:
	"\U0001f4a4",                         // 💤 :zzz:
}
