# imgcompress

Is a simple microservice, that given an image return a link to JPEG image stored
on google cloud storage REGIONAL ( west-ue ). It's ready to be deployed on app 
engine as it is, it uses rpc as a communication method and only exposes one 
rpc call Compress which accepts an image in form of bytes and the quality which
is the compression parameter, from 0 to 100, 100 means as the original 50 means
half of the quality of the original, 50 is a recommended parameter since the 
quality degradation is not "clearly visible to the human eye".
Size reduction the primary the only objective of this service, is not directly
proportional to the quality parameter, i.e an image of 600K is reduced to roughly
100K ( assuming had 100% quality to start with ).

- TODO
	add testing improve the logic/performance;
