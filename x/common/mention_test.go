package common

import (
	"log"
	"testing"
)

func TestCheckMention(t *testing.T) {
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
				target: "init case",
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
				content: " ＃no #Ἀλέξανδρος #KanJAGfånotan #Followme #Facebook #Instagram #Elautobús #Gebäude",
			},
			want: want{
				hashtags: []string{"ἀλέξανδρος", "kanjagfånotan", "followme", "facebook", "instagram", "elautobús", "gebäude"},
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
		{
			in: input{
				target:  "combine tag case",
				content: "@da @1! @2aad~ @#1d21 # @! @ㄎ @} @_ @الجلالة‎",
			},
			want: want{
				hashtags: []string{"1d21"},
				tags:     []string{"da", "1", "2aad", "_"},
			},
		},
		{
			in: input{
				target:  "common case",
				content: "‎tgif#tgif#happy#快樂#hen棒 thanks@zulu111,@kevin994",
			},
			want: want{
				hashtags: []string{"tgif", "happy", "快樂", "hen棒"},
				tags:     []string{"zulu111", "kevin994"},
			},
		},
	}

	for _, test := range tests {
		t.Run("common.CheckMention", func(t *testing.T) {
			hashtags, tags := common.CheckMention(test.in.content)
			log.Println("want", test.want.hashtags, test.want.tags)
			log.Println("out", hashtags, tags)
			if len(hashtags) != len(test.want.hashtags) {
				t.Errorf("len hashtags() = %v, want %v", len(hashtags), len(test.want.hashtags))
			}
			for i := range test.want.hashtags {
				if hashtags[i] != test.want.hashtags[i] {
					t.Errorf("hashtags() = %v, want %v", hashtags, test.want.hashtags)
				}
			}
			if len(tags) != len(test.want.tags) {
				t.Errorf("len tags() = %v, want %v", len(tags), len(test.want.tags))
			}
			for i := range test.want.tags {
				if tags[i] != test.want.tags[i] {
					t.Errorf("tags() = %v, want %v", hashtags, test.want.tags)
				}
			}
		})
	}
}
