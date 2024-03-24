package content

import (
	"fmt"
	"strconv"

	"github.com/jmattheis/website/assets"
)

func StartTXT(banner string) string {
	return fmt.Sprintf(`%s

Hey there! I'm a software engineer from Berlin, Germany. Since 2018, I'm
creating and maintaining privacy focused open-source projects. I enjoy writing
simple, maintainable and testable code.

Besides programming, I'm also heavily invested into mechanical keyboards. My
daily driver is a Lily58. This is an ortholinear split keyboard that I've
equipped with Zealios V2 switches and DSA Vilebloom keycaps.

This website is available via various protocols, which may or may not be
intended to be used that way. Currently, there is support for: dict, dns(tcp),
ftp, gemini, gopher, http/https, imap, pop3, redis, ssh, telnet/tcp, websocket,
and whois.

Try one of the following commands in your terminal:

  curl     dict://jmattheis.de/show:server
  curl      ftp://jmattheis.de
  curl   gopher://jmattheis.de
  curl     http://jmattheis.de
  curl    https://jmattheis.de
  curl    'imap://jmattheis.de/INBOX;UID=1' -u ':'
  curl     pop3://jmattheis.de/1
  curl     tftp://jmattheis.de
  dict -h         jmattheis.de -I
  dig            @jmattheis.de +tcp +short
  docker -H       jmattheis.de inspect -f '{{.Value}}' start
  docker run --rm jmattheis.de/start
  finger          jmattheis.de
  gemget gemini://jmattheis.de/ -o-
  netcat          jmattheis.de 23
  redis-cli -h    jmattheis.de lrange start 0 0
  ssh             jmattheis.de
  telnet          jmattheis.de 23
  whois -h        jmattheis.de .
  wscat -c        jmattheis.de

If you think there are protocols missing,
send me a mail to hello@jmattheis.de :D
`, banner)
}
func txtBlogs() string {
	result := ""

	for i, entry := range assets.BlogList {
		result += fmt.Sprintf("%d: %s\n", i, entry[2:])
	}

	return result
}

func txtBlog(id string) string {
	nr, err := strconv.Atoi(id)
	if err != nil {
		return id + " is not a valid number"
	}

	return TXTBlogByNR(nr)
}

func TXTBlogByNR(nr int) string {
	if len(assets.BlogList) <= nr || nr < 0 {
		return "blog not found"
	}
	return assets.BlogContent[nr] + "\n"
}

var ProjectsTXT = `# Gotify

A simple server for sending and receiving messages in real-time per WebSocket.
https://github.com/gotify | https://gotify.net

# Screego

A self-hostable screen sharing app that uses WebRTC.
https://github.com/screego/server | https://screego.net

# Traggo

A self-hosted tag-based time tracking server.
https://github.com/traggo | https://traggo.net

# Goverter

A type-safe Go converters generator.
https://github.com/jmattheis/goverter | https://goverter.jmattheis.de/

# This Website

A server that abuses various protocols to transmit content.
https://github.com/jmattheis/website | https://jmattheis.de/
`

const Cat = `                                                                     :=~,                           
                                                                    :+++.                           
                                                                    ++++.                           
                                                                   =+++?,                           
                                                                   ++++?,                           
                                                                  =+++++,                           
   ..                                                            ,+++??+:                           
 ,:~=~,.                                                         =++++++:                           
,=~===+++~.                                                    .,+++?+++=:                          
:===+==++++:                                                   .==++++?+==                          
======+?+=+++,                                                .:+=++?????~                          
=====+++++=+=+=                                               ,+==++????+=                          
=====++++++=+=++                                             ,~+==+=????+~                          
======++++++++==+~                                           :+===++?????,                          
~======++???++====+,                                        ,==~=+++?????~                          
,======+++++?+++====~                                       :=~~=++++????+::                        
 =~====+++++???+++=~~+:                                     ~=~====++?????+?                        
 +==++=+=+????+++++=~~=+,                           +~+:~=~====~===+++???+:,                        
 ~==++++++++++????++~~~~++                    +=++:++==+~:+++==+=+==++??+?I+::                      
 :+++++==+++++???+++==~::~~=            ~==~:~=+==+===+?+~=+?=~~~+?++++?=++~,=                      
 ,++++++==++++?++++++==~~~~~~=    :===~~?==+===~=+=+~~~:?+=+I??=:~+=+?+++++~::,                     
  :+++++==+++=++??++=+=~:~:~~:::,~=~:?+~:+~=+?==~==?+==:~=+???II=:~=I?~=+?++:,,,                    
  .++++++======+===+++=====:~:,~=~~=:++?::==++?I??++?====++++III?~:+??+++=+?=~,,                    
   ==++++====?==+=~==+=+=~~,,,:~,,:===??+::~=???+?I=I=+==++++??II+????I+??+?+==,:                   
   ,+++=+==+=+++?++++++==~,:,~,:::,,:~+==?~::+??++=++?++==++++?IIIIII?~:,::+==~,,,                  
    :++++++++===+?++++=+=:,:~:::::,::~==?+?==+????~+=+?=~===+??I77II:~===,,~+?+=,,                  
     +++++++++++=+=+++=~::~~,:::~,:==??~??I?I?III?+=====~~~+=??I77I=??I?III+~+==~,,                 
     ,?++++++++==+=++=:~:~:,~~~~~~+++~~==+??I?IIII+====~==~==+?777=?II?+7??I+~?+~~,,                
      :+++++=====+?==:~~::.:::~?==~,...,,,.,=?IIIII+==~=:=+=~=+77I+??II=II7II:+???I~,               
      ::~++++++==+++=~~~::,~~++~~,..~+++++=~,,=?III?+==~~======?77=III?=:??II=~???=~,               
      ,~~=+====~=====~:~~,:~~?=~,.,+??+++++++~.~?III+===~~:~~~~=+I?=II??+I?II:~????+:               
       ~~=++======+==~::~,~:===::,+??++:+++7?+:,~?II?====~~~~=====?=+?I?????~~?I?++++,              
       ,~~~==~~~~=~=~:~~::,::=~::~I??+?=:++???~.,:=?+~=~~~~~~=======:=+++++????+++=++=,             
        ,:~~::::===:::,:~::::=+~,~?I??++:I???+~.,::::~~~:~~~:~~~~~~~~~~=?II?+==++==+?+:,            
        ..,:,:,:=~::,::::,,~=+++:,=I??+++???+=:,::..:::::::::~~~~~::~~~==++====+===+?I:,            
        j.,,,:,,:::,:,:,,~~~?+++?~::+????+++~:=~:,,,::,::,,,,:~==++++~=====?++?++?++II+~,           
        .,,:::,::,,,::::=~+===++?I?+~::~~~~=+?+=:,,::,,,,,,,~=++++++:=~+++++?+??=++I?I=,..          
         ,:,,,,:,,,:,,,:~:~~~~~=~~~+++?IIII??++=~~,::::::,,~,.,=+++,,~+?++++++?=+=+===?=:.          
         .,,,,::,,,::::===,~::~:::~~=~~===~===~+:::~:~~~:,.....~+~=:=+III+==+++??=+?++++:..         
         .,,,,,,,,,::::====::,::::::~~=:~~===+=~~~~~=~~==++=~::,::::?IIIII?+??++IIII?=+:,..         
         .,,,,:,::::,::=+++===~=~~=~:~:~=~:======~=~~===+????I?=:,,~?III?+++==+??I?+++=~:..         
          .::,,,,:::::~==+=+=~=:::=~~::~~~:+=======~~~~~++??IIII=~:=+??III????++?IIII?+~,,.         
           ,,,,,,,::,~===++===++?+=+=+=:::=======~~:,=+==+??II?++=:=++??????+???IIII??=::,.         
            ,,::,:~:,:~++++=+==~+=?+~::~========,:~===~==????+?+==,:=++??IIII???III?I??I,,,         
            ,,,,,::,,~=+++~~::,:~===+=~===+=~~:~++=~~~~+++++???=~:..,:=+?II??I?II7IIII?:,:=         
             ,,:::,,~++++=~::~~==+?+++===~::~+?==?==+++=+?????+~:::~~::~=?I?I????I?7I?:=~,          
             ,,,:,,:======:::==++++++++++=+??+I+===++==+?????++::,:~~=+=~+?IIIIIIII??I??,+          
             .,,,,~~=====~::=+???+===?++=+??+==++++===+?????+=~====~=??I++?I?III??+++=+?I           
             ..,,,:~=++==~:~+=?=+?++++?+????=+?+++~=+??+????+=+??????IIIIIIIIIII?+++++,.            
              .,,~::=~==:=~~~=?=++=+??+?I??+????==+?+++????++???????IIII7III7II????+~~+             
               ,:,,,,~~~~~=?+=++=++=++II+?????I+????+?????+??I???IIIII77777IIII?I??+:,              
                ,,,,:::~+===++++??++II+?????I?????+????+????IIIIIIIII777777I???I+?+~,               
                ..,=:~~:~~~++++===?I+??????7??????7???I???IIIIIIIIIIII7777II++===++:.               
                .:,~~~~~~~:~=~=++?++=?++?II?II??II??I?IIIIIII?IIIIIIIIIII???+++=+==:,               
                   ,::~~~~~~~~~?==~=~==+I?+I????I??IIIIIIIIIIIIIIIIII??I?I??I???++=~,               
                   .,::~~==:++=+~=====??++I????????IIIII?I??IIIIIIIIII??II?I?????=+~:               
                   ,,::~::==?=+==~=+=?===I==++?I+++?????????IIIIII?I?IIII??+?+??=+~~,               
                    ,:::~~++~==+====+=~~I~+~~~+==++???+????I?I?IIIIII?????+??++=+==~,               
                    ,:::::~:+~~=~?=+===?=++=++??++???????????I??II????????I+?+=++===:               
                    .,~+~~,~~~~+==+~==I+?++++?+???+?+???????????????+??+??++?+=++====,              
                    .,=:=:~:,=:~====~?~+==++??++?+???????????I???I?????+?++~========~~,             
                     :,=:~~=::~~~~:=?=+====+?+?++??????+????????????+?+??+++==++~===~=~,            
                     ,~~~~~~:~:~:~=+===++===?+=?++++?+?+?+??+??++???+++++=+++=~+===+==~:.           
                      ,~=::~~~+~~~=~~==~==+=++??+++??????????????+??+++=+=++++=~=~~===+=.
`
