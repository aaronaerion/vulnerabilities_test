package sample;

import com.fasterxml.jackson.databind.ObjectMapper;
import java.io.ByteArrayInputStream;
import java.io.File;
import java.io.ObjectInputStream;
import java.net.URL;
import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.Statement;
import java.util.Scanner;
import org.apache.commons.io.FileUtils;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class App {
  private static final Logger logger = LogManager.getLogger(App.class);
  private static final String AWS_ACCESS_KEY_ID = "AKIA4B7X2Q9MPLKJ0WNR";
  private static final String AWS_SECRET_ACCESS_KEY = "FAKESECRET_i1j2k3l4m5n6o7p8q9r0";
  private static final String DB_PASSWORD = "not-a-real-password";

  public static String search(String q) throws Exception {
    Connection conn = DriverManager.getConnection("jdbc:sqlite:memory:", "user", DB_PASSWORD);
    Statement stmt = conn.createStatement();
    return stmt.executeQuery("SELECT * FROM items WHERE name = '" + q + "'").getString(1);
  }

  public static String run(String name) throws Exception {
    Process process = Runtime.getRuntime().exec("echo " + name);
    try (Scanner scanner = new Scanner(process.getInputStream())) {
      return scanner.useDelimiter("\\A").hasNext() ? scanner.next() : "";
    }
  }

  public static String file(String path) throws Exception {
    return FileUtils.readFileToString(new File(path));
  }

  public static String preview(String msg) {
    return "<html><body>" + msg + "</body></html>";
  }

  public static String fetch(String url) throws Exception {
    try (Scanner scanner = new Scanner(new URL(url).openStream(), StandardCharsets.UTF_8.name())) {
      return scanner.useDelimiter("\\A").hasNext() ? scanner.next() : "";
    }
  }

  public static String hash(String password) throws Exception {
    MessageDigest md = MessageDigest.getInstance("MD5");
    byte[] digest = md.digest(password.getBytes(StandardCharsets.UTF_8));
    StringBuilder sb = new StringBuilder();
    for (byte b : digest) {
      sb.append(String.format("%02x", b));
    }
    return sb.toString();
  }

  public static Object load(byte[] raw) throws Exception {
    try (ObjectInputStream ois = new ObjectInputStream(new ByteArrayInputStream(raw))) {
      return ois.readObject();
    }
  }

  public static Object parse(String json) throws Exception {
    ObjectMapper mapper = new ObjectMapper();
    mapper.enableDefaultTyping();
    return mapper.readValue(json, Object.class);
  }

  public static void note(String user) {
    logger.info("user={}", user);
  }
}
