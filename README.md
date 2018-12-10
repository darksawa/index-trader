# Index Trader

Figures out at what date and at what price to buy and sell the index in order to get the highest return. Only one buy and one sell transaction is allowed in that order (short selling is not allowed).

Assume that:
- intraday trading is not possible, otherwise we cannot guarantee that it is not a short sell, because candlestick data is not enough for that
- we can buy on the lost and sell on highest price :)

Algorithm:
- Initialization: we buy on the lowest price of day 1, sell on the highest price of day 2 and mark this deal as the best and current 
- starting from day 3, check the lowest price of the previous day, 
 if we found new lowest price - compare current deal with the best one, if it is better, mark the current one as the best,
 close the current deal and open a new one
- if the highest price of current date is higher than in the current deal's sell price, sell on that price )
- increment day counter.
   

Time complexity O(n)