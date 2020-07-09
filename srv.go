package main

/*
client ->:
	impulse index						Player choice

server ->:
	join			player name piece	Player joins
	leave			player				Player leaves
	turn			player				Turn begins
	setown			square player		Square changes hands
	setupgrade		square n			Square house/hotel level is set to n
	addmortgage		square				Square gets mortgaged
	auction			square				Auction begins
	bid				player amount		Player bids in the auction
*/
