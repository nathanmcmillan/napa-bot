
import java.io.BufferedReader;
import java.io.DataOutputStream;
import java.io.InputStreamReader;
import java.net.URL;
import javax.net.ssl.HttpsURLConnection;

import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;
import java.util.Base64;

public class Gdax {
    private static final String apiRest = "https://api.gdax.com";
    private static final String apiFeed = "wss://ws-feed.gdax.com";
    private static final String sandRest = "https://public.sandbox.gdax.com";
    private static final String sandFeed = "wss://ws-feed-public.sandbox.gdax.com";
    
    private static final String apiCurrencies = "/currencies";
    private static final String apiProducts = "/products";
    private static final String apiTrades = "/trades";
    private static final String apiFills = "/fills";
    private static final String apiOrders = "/orders";
    private static final String apiAccounts = "/accounts";
    private static final String apiTime = "/time";
    
    private static final String accessKey = "CB-ACCESS-KEY";
    private static final String accessSign = "CB-ACCESS-SIGN";
    private static final String accessTime = "CB-ACCESS-TIMESTAMP";
    private static final String accessPass = "CB-ACCESS-PASSPHRASE";
    
    public static final String btcUsd = "BTC-USD";
    public static final String ethUsd = "ETH-USD";
    public static final String ltcUsd = "LTC-USD";

    private static final String get = "GET";
    private static final String post = "POST";
    
    public static final String onlyBest = "1";
    public static final String top50 = "2";
    public static final String fullBook = "3";
    
    public static final int codeSuccess = 200;
    public static final int codeBadRequest = 400;
    public static final int codeUnauthorized = 401;
    public static final int codeForbidden = 403;
    public static final int codeNotFound = 404;
    public static final int codeRateLimit = 429;
    public static final int codeServerError = 500;
    
    public void getProducts() throws Exception {
        publicSend(get, apiRest + apiProducts);
    }
    
    public void getProductOrderBook(String productId, String level) throws Exception {
        StringBuilder s = new StringBuilder();
        s.append(apiRest).append(apiProducts).append("/").append(productId).append("/book");
        if (level != null) {
            s.append("?level=").append(level);
        }
        publicSend(get, s.toString());
    }
    
    public void getProductTicker(String productId) throws Exception {
        publicSend(get, apiRest + apiProducts + "/" + productId + "/ticker");
    }
    
    public void getProductTrades(String productId) throws Exception {
        publicSend(get, apiRest + apiProducts + "/" + productId + "/trades");
    }
    
    public void getProductStats(String productId) throws Exception {
        publicSend(get, apiRest + apiProducts + "/" + productId + "/stats");
    }
    
    public void getTime() throws Exception {
        publicSend(get, apiRest + apiTime);
    }
    
    public void getCurrencies() throws Exception {
        publicSend(get, apiRest + apiCurrencies);
    }
    
    public void getOrders(String key, String sign, String time, String pass) throws Exception {
        privateSend(get, apiRest + apiOrders, key, sign, time, pass);
    }
    
    public void getTrades(String key, String sign, String time, String pass) throws Exception {
        privateSend(get, apiRest + apiTrades, key, sign, time, pass);
    }
    
    public void getFills(String key, String sign, String time, String pass) throws Exception {
        privateSend(get, apiRest + apiFills, key, sign, time, pass);
    }
    
    public void getAccounts(String key, String sign, String time, String pass) throws Exception {
        privateSend(get, apiRest + apiAccounts, key, sign, time, pass);
    }
    
    private static void privateSend(String method, String address, String key, String sign, String time, String pass) throws Exception {
        URL url = new URL(address);
        HttpsURLConnection connection = (HttpsURLConnection)url.openConnection();
        connection.setRequestMethod(method);
        connection.setRequestProperty(accessKey, key);
        connection.setRequestProperty(accessSign, sign);
        connection.setRequestProperty(accessTime, time);
        connection.setRequestProperty(accessPass, pass);
        int responseCode = connection.getResponseCode();
        System.out.println(method + " " + url);
        System.out.println("code " + responseCode);
        try (BufferedReader in = new BufferedReader(new InputStreamReader(connection.getInputStream()))) {
            String inputLine = in.readLine();
            StringBuffer response = new StringBuffer();
            while (inputLine != null) {
                response.append(inputLine);
                inputLine = in.readLine();
            }
            System.out.println("message " + response.toString());
        }
    }
    
    private static void publicSend(String method, String address) throws Exception {
        URL url = new URL(address);
        HttpsURLConnection connection = (HttpsURLConnection)url.openConnection();
        connection.setRequestMethod(method);
        int responseCode = connection.getResponseCode();
        System.out.println(method + " " + url);
        System.out.println("code " + responseCode);
        try (BufferedReader in = new BufferedReader(new InputStreamReader(connection.getInputStream()))) {
            String inputLine = in.readLine();
            StringBuffer response = new StringBuffer();
            while (inputLine != null) {
                response.append(inputLine);
                inputLine = in.readLine();
            }
            System.out.println("message " + response.toString());
        }
    }
    
    private static String signature() {
        
    }
}
