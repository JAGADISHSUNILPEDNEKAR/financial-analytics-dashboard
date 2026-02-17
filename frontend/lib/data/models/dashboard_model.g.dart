// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'dashboard_model.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_$DashboardModelImpl _$$DashboardModelImplFromJson(Map<String, dynamic> json) =>
    _$DashboardModelImpl(
      id: json['id'] as String,
      userId: json['userId'] as String,
      name: json['name'] as String,
      layout: json['layout'] as Map<String, dynamic>,
      widgets: (json['widgets'] as List<dynamic>)
          .map((e) => WidgetModel.fromJson(e as Map<String, dynamic>))
          .toList(),
      isPublic: json['isPublic'] as bool,
      createdAt: DateTime.parse(json['createdAt'] as String),
      updatedAt: DateTime.parse(json['updatedAt'] as String),
    );

Map<String, dynamic> _$$DashboardModelImplToJson(
  _$DashboardModelImpl instance,
) =>
    <String, dynamic>{
      'id': instance.id,
      'userId': instance.userId,
      'name': instance.name,
      'layout': instance.layout,
      'widgets': instance.widgets,
      'isPublic': instance.isPublic,
      'createdAt': instance.createdAt.toIso8601String(),
      'updatedAt': instance.updatedAt.toIso8601String(),
    };

_$WidgetModelImpl _$$WidgetModelImplFromJson(Map<String, dynamic> json) =>
    _$WidgetModelImpl(
      id: json['id'] as String,
      type: json['type'] as String,
      config: json['config'] as Map<String, dynamic>,
      position: json['position'] as Map<String, dynamic>,
    );

Map<String, dynamic> _$$WidgetModelImplToJson(_$WidgetModelImpl instance) =>
    <String, dynamic>{
      'id': instance.id,
      'type': instance.type,
      'config': instance.config,
      'position': instance.position,
    };
