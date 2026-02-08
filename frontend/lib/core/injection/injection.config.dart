// GENERATED CODE - DO NOT MODIFY BY HAND

// **************************************************************************
// InjectableConfigGenerator
// **************************************************************************

// ignore_for_file: type=lint
// coverage:ignore-file

// ignore_for_file: no_leading_underscores_for_library_prefixes
import 'package:get_it/get_it.dart' as _i174;
import 'package:injectable/injectable.dart' as _i526;

import '../../presentation/blocs/auth/auth_bloc.dart' as _i141;
import '../../presentation/blocs/connectivity/connectivity_bloc.dart' as _i905;
import '../../services/auth_service.dart' as _i610;

extension GetItInjectableX on _i174.GetIt {
// initializes the registration of main-scope dependencies inside of GetIt
  _i174.GetIt init({
    String? environment,
    _i526.EnvironmentFilter? environmentFilter,
  }) {
    final gh = _i526.GetItHelper(
      this,
      environment,
      environmentFilter,
    );
    gh.factory<_i905.ConnectivityBloc>(() => _i905.ConnectivityBloc());
    gh.lazySingleton<_i610.AuthService>(() => _i610.AuthServiceImpl());
    gh.factory<_i141.AuthBloc>(() => _i141.AuthBloc(gh<_i610.AuthService>()));
    return this;
  }
}
