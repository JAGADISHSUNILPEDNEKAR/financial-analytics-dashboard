import 'package:injectable/injectable.dart';

abstract class AuthService {
  Future<bool> get isAuthenticated;
}

@LazySingleton(as: AuthService)
class AuthServiceImpl implements AuthService {
  @override
  Future<bool> get isAuthenticated async {
    // TODO: Implement actual auth check
    return false;
  }
}
