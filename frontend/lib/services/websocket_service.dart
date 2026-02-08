import 'dart:async';
import 'dart:convert';
import 'package:web_socket_channel/web_socket_channel.dart';
import 'package:financial_analytics/core/constants/api_constants.dart';
import 'package:financial_analytics/services/auth_service.dart';
import 'package:get_it/get_it.dart';

class WebSocketService {
  static final WebSocketService _instance = WebSocketService._internal();
  factory WebSocketService() => _instance;
  WebSocketService._internal();

  WebSocketChannel? _channel;
  final _messageController = StreamController<Map<String, dynamic>>.broadcast();
  final _connectionController = StreamController<ConnectionStatus>.broadcast();
  Timer? _reconnectTimer;
  Timer? _pingTimer;
  bool _intentionalClose = false;
  int _reconnectAttempts = 0;
  
  Stream<Map<String, dynamic>> get messages => _messageController.stream;
  Stream<ConnectionStatus> get connectionStatus => _connectionController.stream;
  
  Future<void> connect() async {
    if (_channel != null) return;
    
    try {
      final token = await GetIt.instance<AuthService>().getAccessToken();
      if (token == null) throw Exception('No auth token available');
      
      final wsUrl = Uri.parse('${ApiConstants.wsBaseUrl}/ws?token=$token');
      _channel = WebSocketChannel.connect(wsUrl);
      
      _connectionController.add(ConnectionStatus.connecting);
      
      _channel!.stream.listen(
        _handleMessage,
        onDone: _handleDisconnect,
        onError: _handleError,
      );
      
      _connectionController.add(ConnectionStatus.connected);
      _reconnectAttempts = 0;
      _startPingTimer();
      
    } catch (e) {
      _connectionController.add(ConnectionStatus.disconnected);
      _scheduleReconnect();
    }
  }
  
  void _handleMessage(dynamic message) {
    try {
      final data = json.decode(message as String);
      _messageController.add(data);
    } catch (e) {
      print('Error parsing WebSocket message: $e');
    }
  }
  
  void _handleDisconnect() {
    _channel = null;
    _connectionController.add(ConnectionStatus.disconnected);
    _stopPingTimer();
    
    if (!_intentionalClose) {
      _scheduleReconnect();
    }
  }
  
  void _handleError(error) {
    print('WebSocket error: $error');
    _connectionController.add(ConnectionStatus.error);
    disconnect();
  }
  
  void _scheduleReconnect() {
    if (_reconnectAttempts >= 5) return;
    
    final delay = Duration(seconds: 2 * (_reconnectAttempts + 1));
    _reconnectTimer = Timer(delay, () {
      _reconnectAttempts++;
      connect();
    });
  }
  
  void _startPingTimer() {
    _pingTimer = Timer.periodic(const Duration(seconds: 30), (_) {
      send({'type': 'ping'});
    });
  }
  
  void _stopPingTimer() {
    _pingTimer?.cancel();
    _pingTimer = null;
  }
  
  void send(Map<String, dynamic> message) {
    if (_channel != null) {
      _channel!.sink.add(json.encode(message));
    }
  }
  
  void subscribe(List<String> symbols) {
    send({
      'type': 'subscribe',
      'symbols': symbols,
    });
  }
  
  void unsubscribe(List<String> symbols) {
    send({
      'type': 'unsubscribe',
      'symbols': symbols,
    });
  }
  
  void disconnect() {
    _intentionalClose = true;
    _reconnectTimer?.cancel();
    _stopPingTimer();
    _channel?.sink.close();
    _channel = null;
    _connectionController.add(ConnectionStatus.disconnected);
  }
  
  void dispose() {
    disconnect();
    _messageController.close();
    _connectionController.close();
  }
}

enum ConnectionStatus { disconnected, connecting, connected, error }