package main

import (
	"log"
	"testing"
)

func Test_checkMention(t *testing.T) {
	type input struct {
		target  string
		content string
	}
	type want struct {
		hashtags []string
		tags     []string
	}
	tests := []struct {
		in   input
		want want
	}{
		{
			in: input{
				target: "common case",
				content: "你好123 11dDD	#413123, #_1#r1 #😘 #⑽ # #⓱ #◭ #🉑 #。 #! #_ #666",
			},
			want: want{
				hashtags: []string{"413123", "_1", "r1", "😘", "⑽", "⓱", "◭", "🉑", "_", "666"},
				tags:     []string{},
			},
		},
		{
			in: input{
				target:  "asia case",
				content: "＃にほんご #本屋で売る小説#1本書  #靠谱~ a#彼が食べない寿司 #감사합니다 #শ螈襎胂ㄝ #ตัวเอง, #ᠮᠣᠩᠭᠣᠯ",
			},
			want: want{
				hashtags: []string{"本屋で売る小説", "1本書", "靠谱", "彼が食べない寿司", "감사합니다", "শ螈襎胂ㄝ", "ตัวเอง", "ᠮᠣᠩᠭᠣᠯ"},
				tags:     []string{},
			},
		},
		{
			in: input{
				target:  "west case",
				content: " ＃no #Ἀλέξανδρος #KanJAGfånotan #Facebook #Instagram #Elautobús #Gebäude",
			},
			want: want{
				hashtags: []string{"ἀλέξανδρος", "kanjagfånotan", "facebook", "instagram", "elautobús", "gebäude"},
				tags:     []string{},
			},
		},
		{ // Punctuation ref: https://en.wikipedia.org/wiki/Punctuation ("ฯ", "º", "ª" are letters)
			in: input{
				target: "punctuation",
				content: `#’ # #' #[ #] #( #) #{ #} #⟨ #⟩
				#: #, #، #、 #‒ #– #— #― #… #... #⋯ #᠁ #! #. #‹ #› #« #» #‐ #- #?
				#‘ #’ #“ #” #' #' #" #" #; #/ #⧸ #⁄ #· #& #* #@ #\ #‱ #• #^ #† #‡ #⹋ #°
				#” #= #¡ #¿ #※ #× ## #№ #÷ #% #‰ #+ #− #± #∓ #′ #″ #‴ #§ #~ #¶ #| #‖ #¦`,
			},
			want: want{
				hashtags: []string{},
				tags:     []string{},
			},
		},
	}

	for _, test := range tests {
		t.Run("checkMention", func(t *testing.T) {
			hashtags, _ := checkMention(test.in.content)
			log.Println("want", test.want.hashtags)
			log.Println("out", hashtags)
			if len(hashtags) != len(test.want.hashtags) {
				t.Errorf("len hashtags() = %v, want %v", len(hashtags), len(test.want.hashtags))
			}
			for i := range test.want.hashtags {
				if hashtags[i] != test.want.hashtags[i] {
					t.Errorf("hashtags() = %v, want %v", hashtags, test.want.hashtags)
				}
			}
		})
	}
}
