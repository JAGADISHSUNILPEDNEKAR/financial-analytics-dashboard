import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:injectable/injectable.dart';

abstract class ConnectivityEvent {}
class ConnectivityCheckRequested extends ConnectivityEvent {}

abstract class ConnectivityState {}
class ConnectivityInitial extends ConnectivityState {}

@injectable
class ConnectivityBloc extends Bloc<ConnectivityEvent, ConnectivityState> {
  ConnectivityBloc() : super(ConnectivityInitial()) {
    on<ConnectivityCheckRequested>((event, emit) {
      // TODO: Implement connectivity check
    });
  }
}
