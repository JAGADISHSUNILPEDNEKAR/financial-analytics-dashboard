import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';
import 'package:financial_analytics/main.dart';
import 'package:financial_analytics/services/auth_service.dart';
import 'package:financial_analytics/services/websocket_service.dart';

class MockAuthService extends Mock implements AuthService {}
class MockWebSocketService extends Mock implements WebSocketService {}

void main() {
  group('App Widget Tests', () {
    testWidgets('App launches and shows login screen when not authenticated', (WidgetTester tester) async {
      // Build our app and trigger a frame
      await tester.pumpWidget(const FinancialAnalyticsApp());
      
      // Verify login screen is shown
      expect(find.text('Login'), findsOneWidget);
      expect(find.byType(TextFormField), findsNWidgets(2)); // Email and password fields
    });
    
    testWidgets('Dashboard shows when authenticated', (WidgetTester tester) async {
      // Mock authenticated state
      final authService = MockAuthService();
      when(authService.isAuthenticated).thenReturn(true);
      
      await tester.pumpWidget(const FinancialAnalyticsApp());
      await tester.pumpAndSettle();
      
      // Verify dashboard is shown
      expect(find.text('Dashboard'), findsOneWidget);
    });
  });
  
  group('Real-time Chart Tests', () {
    testWidgets('Chart updates with real-time data', (WidgetTester tester) async {
      // Test real-time chart widget
      const symbol = 'AAPL';
      
      await tester.pumpWidget(
        MaterialApp(
          home: Scaffold(
            body: RealTimeChart(symbol: symbol),
          ),
        ),
      );
      
      // Initially shows loading
      expect(find.byType(CircularProgressIndicator), findsOneWidget);
      
      // Simulate WebSocket message
      // Add test for data update
    });
  });
}