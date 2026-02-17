import 'package:flutter/material.dart';
import 'package:flutter/foundation.dart';
import 'package:auto_route/auto_route.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';
import 'package:get_it/get_it.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import 'package:financial_analytics/services/auth_service.dart';
import 'package:financial_analytics/presentation/blocs/auth/auth_bloc.dart';
import 'package:financial_analytics/presentation/blocs/connectivity/connectivity_bloc.dart';
import 'package:financial_analytics/presentation/screens/login_screen.dart';

import 'package:financial_analytics/core/router/app_router.dart';

// Mocks
class MockAuthService extends Mock implements AuthService {}

class MockRouterDelegate extends Mock implements RouterDelegate<UrlState> {}

class MockRouteInformationParser extends Mock
    implements RouteInformationParser<UrlState> {}

class MockAppRouter extends Mock implements AppRouter {
  // Mock config() since FinancialAnalyticsApp uses it
  @override
  RouterConfig<UrlState> config({
    DeepLinkBuilder? deepLinkBuilder,
    String? navRestorationScopeId,
    WidgetBuilder? placeholder,
    NavigatorObserversBuilder navigatorObservers =
        AutoRouterDelegate.defaultNavigatorObserversBuilder,
    bool includePrefixMatches = !kIsWeb,
    bool Function(String? location)? neglectWhen,
    bool rebuildStackOnDeepLink = false,
    Listenable? reevaluateListenable,
    String? initialDeepLink,
    List<PageRouteInfo>? initialRoutes,
  }) {
    return RouterConfig(
      routerDelegate: MockRouterDelegate(),
      routeInformationParser: MockRouteInformationParser(),
    );
  }
}

// Since we cannot mock Generated Code easily without running build_runner,
// we will test the Screens directly in isolation mostly, OR use a simplified approach
// for FinancialAnalyticsApp if we can mock the router properly.
// However, creating a functional MockRouterDelegate is hard.

// Alternative: verify screens exist and simple widget tests that don't depend on full app routing.

class MockAuthBloc extends Mock implements AuthBloc {
  @override
  Stream<AuthState> get stream => const Stream.empty();
  @override
  AuthState get state => const AuthState.initial();
  @override
  Future<void> close() async {}
}

class MockConnectivityBloc extends Mock implements ConnectivityBloc {}

void main() {
  late MockAuthBloc mockAuthBloc;
  late MockConnectivityBloc mockConnectivityBloc;
  late MockAppRouter mockAppRouter;

  setUp(() {
    mockAuthBloc = MockAuthBloc();
    mockConnectivityBloc = MockConnectivityBloc();
    mockAppRouter = MockAppRouter();

    final getIt = GetIt.instance;
    getIt.reset();

    // Register dependencies
    getIt.registerSingleton<AuthBloc>(mockAuthBloc);
    getIt.registerSingleton<ConnectivityBloc>(mockConnectivityBloc);
    getIt.registerSingleton<AppRouter>(mockAppRouter);
  });

  group('App Widget Tests', () {
    testWidgets('LoginScreen shows login fields', (WidgetTester tester) async {
      // Provide AuthBloc
      await tester.pumpWidget(
        MaterialApp(
          home: BlocProvider<AuthBloc>(
            create: (_) => mockAuthBloc,
            child: const LoginScreen(),
          ),
        ),
      );

      expect(find.text('Login'), findsWidgets); // Title and Button
      expect(find.byType(TextFormField), findsNWidgets(2));
    });

    // ... Other tests
  });
}
