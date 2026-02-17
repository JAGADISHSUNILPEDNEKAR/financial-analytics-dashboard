import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:financial_analytics/domain/entities/dashboard.dart';

part 'dashboard_model.freezed.dart';
part 'dashboard_model.g.dart';

@freezed
class DashboardModel with _$DashboardModel {
  const DashboardModel._();

  const factory DashboardModel({
    required String id,
    required String userId,
    required String name,
    required Map<String, dynamic> layout,
    required List<WidgetModel> widgets,
    required bool isPublic,
    required DateTime createdAt,
    required DateTime updatedAt,
  }) = _DashboardModel;

  factory DashboardModel.fromJson(Map<String, dynamic> json) =>
      _$DashboardModelFromJson(json);

  Dashboard toEntity() => Dashboard(
        id: id,
        userId: userId,
        name: name,
        layout: layout,
        widgets: widgets.map((w) => w.toEntity()).toList(),
        isPublic: isPublic,
        createdAt: createdAt,
        updatedAt: updatedAt,
      );
}

@freezed
class WidgetModel with _$WidgetModel {
  const WidgetModel._();

  const factory WidgetModel({
    required String id,
    required String type,
    required Map<String, dynamic> config,
    required Map<String, dynamic> position,
  }) = _WidgetModel;

  factory WidgetModel.fromJson(Map<String, dynamic> json) =>
      _$WidgetModelFromJson(json);

  Widget toEntity() =>
      Widget(id: id, type: type, config: config, position: position);
}
