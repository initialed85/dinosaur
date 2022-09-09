import java.io.IOException;
import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetSocketAddress;
import java.net.SocketAddress;
import java.nio.charset.StandardCharsets;

class App {
    final int PORT = 13337;

    App() throws IOException, InterruptedException {
        String hostname = System.getenv("HOSTNAME");
        String localIp = System.getenv("LOCAL_IP");
        String broadcastIp = System.getenv("BROADCAST_IP");

        DatagramSocket datagramSocket = new DatagramSocket(PORT);
        datagramSocket.setBroadcast(true);

        byte[] buf = String.format(
                "Hello world from Java @ %s",
                hostname
        ).getBytes(StandardCharsets.UTF_8);

        SocketAddress broadcastAddress = new InetSocketAddress(broadcastIp, this.PORT);
        DatagramPacket broadcastPacket = new DatagramPacket(buf, buf.length, broadcastAddress);

        new Thread(() -> {
            try {
                this.receiveLoop(datagramSocket, localIp);
            } catch (IOException e) {
                throw new RuntimeException(e);
            }
        }).start();

        for (; ; ) {
            datagramSocket.send(broadcastPacket);
            Thread.sleep(1000);
        }
    }

    private void receiveCallback(DatagramPacket receivedPacket, String localIp) {
        String ip = receivedPacket.getAddress().getHostAddress();
        int port = receivedPacket.getPort();

        if (ip.equals(localIp) && port == PORT) {
            return;
        }

        String data = new String(receivedPacket.getData(), StandardCharsets.UTF_8);

        System.out.format("%s:%d\t%s\n", ip, port, data);
    }

    private void receiveLoop(DatagramSocket socket, String localIp) throws IOException {
        for (; ; ) {
            byte[] buf = new byte[65507];
            DatagramPacket receivedPacket = new DatagramPacket(buf, 65507);
            socket.receive(receivedPacket);

            this.receiveCallback(receivedPacket, localIp);
        }
    }
}

public class Main {
    public static void main(String[] args) throws IOException, InterruptedException {
        App app = new App();
    }
}
