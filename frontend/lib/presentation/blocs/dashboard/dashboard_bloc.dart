import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:freezed_annotation/freezed_annotation.dart';

part 'dashboard_bloc.freezed.dart';

@freezed
class DashboardState with _$DashboardState {
  const factory DashboardState.initial() = _Initial;
  const factory DashboardState.loading() = _Loading;
  const factory DashboardState.loaded(Dashboard dashboard) = _Loaded;
  const factory DashboardState.error(String message) = _Error;
}

@freezed
class DashboardEvent with _$DashboardEvent {
  const factory DashboardEvent.loadRequested() = DashboardLoadRequested;
  const factory DashboardEvent.widgetRemoved({required String widgetId}) = DashboardWidgetRemoved;
}

class Dashboard {
  final List<dynamic> widgets;
  Dashboard({this.widgets = const []});
}

class DashboardBloc extends Bloc<DashboardEvent, DashboardState> {
  DashboardBloc() : super(const DashboardState.initial()) {
    on<DashboardLoadRequested>((event, emit) async {
       emit(const DashboardState.loading());
       // Simulate load
       emit(DashboardState.loaded(Dashboard()));
    });
  }
}
