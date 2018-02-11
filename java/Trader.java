
public class Trader {
    public static void main(String[] args) throws Exception {
        final Gdax gdax = new Gdax();
        final String key = "key";
        final String sign = "sign";
        final String time = "utc epoch here";
        final String pass = "pass";
        gdax.getAccounts(key, sign, time, pass);
    }
}
