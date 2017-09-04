namespace go com.game.trade.addrtx

struct GetAddrMsg{
    1: required string coinType;
    2: required i64 uid;
}
struct GetTXMsg{
    1: required string coinType;
    2: required i64 fromUID;
    3: required i64 fromAmount;
    4: required i64 toUID;
    5: required i64 toAmount;
}

service AddrTXService{
    string GetAddr(1: GetAddrMsg msg);
    string GetTX(1: GetTXMsg msg);
}