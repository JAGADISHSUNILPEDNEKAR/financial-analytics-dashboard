import 'package:injectable/injectable.dart';

abstract class AuthService {
  Future<bool> get isAuthenticated;
  Future<String?> getAccessToken();
}

@LazySingleton(as: AuthService)
class AuthServiceImpl implements AuthService {
  @override
  Future<bool> get isAuthenticated async {
    // TODO: Implement actual auth check
    return false;
  }

  @override
  Future<String?> getAccessToken() async {
    // TODO: Implement actual token retrieval
    return 'dummy_token';
  }
}
