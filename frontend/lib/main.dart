import 'package:flutter/material.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:firebase_core/firebase_core.dart';
import 'package:financial_analytics/core/injection/injection.dart';
import 'package:financial_analytics/core/theme/app_theme.dart';
import 'package:financial_analytics/core/router/app_router.dart';
import 'package:financial_analytics/presentation/blocs/auth/auth_bloc.dart';
import 'package:financial_analytics/presentation/blocs/connectivity/connectivity_bloc.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  // Initialize platform-specific services
  await _initializePlatform();

  // Initialize dependency injection
  await configureDependencies();

  // Initialize local storage
  await Hive.initFlutter();

  // Initialize Firebase
  await Firebase.initializeApp();

  runApp(const FinancialAnalyticsApp());
}

Future<void> _initializePlatform() async {
  if (kIsWeb) {
    // Web-specific initialization
  } else if (defaultTargetPlatform == TargetPlatform.windows ||
      defaultTargetPlatform == TargetPlatform.linux ||
      defaultTargetPlatform == TargetPlatform.macOS) {
    // Desktop-specific initialization
  } else {
    // Mobile-specific initialization
  }
}

class FinancialAnalyticsApp extends StatelessWidget {
  const FinancialAnalyticsApp({super.key});

  @override
  Widget build(BuildContext context) {
    final appRouter = getIt<AppRouter>();

    return MultiBlocProvider(
      providers: [
        BlocProvider(
          create: (_) => getIt<AuthBloc>()..add(const AuthCheckRequested()),
        ),
        BlocProvider(create: (_) => getIt<ConnectivityBloc>()),
      ],
      child: MaterialApp.router(
        title: 'Financial Analytics',
        theme: AppTheme.lightTheme,
        darkTheme: AppTheme.darkTheme,
        routerConfig: appRouter.config(),
        debugShowCheckedModeBanner: false,
      ),
    );
  }
}
