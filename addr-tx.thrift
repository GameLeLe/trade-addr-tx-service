namespace go com.game.trade.addrtx

struct GetAddrMsg{
    1: required string coinType;
    2: required i64 uid;
}

service AddrTXService{
    string GetAddr(1: GetAddrMsg msg);
}