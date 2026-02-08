// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'dashboard_model.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

T _$identity<T>(T value) => value;

final _privateConstructorUsedError = UnsupportedError(
    'It seems like you constructed your class using `MyClass._()`. This constructor is only meant to be used by freezed and you are not supposed to need it nor use it.\nPlease check the documentation here for more information: https://github.com/rrousselGit/freezed#adding-getters-and-methods-to-our-models');

DashboardModel _$DashboardModelFromJson(Map<String, dynamic> json) {
  return _DashboardModel.fromJson(json);
}

/// @nodoc
mixin _$DashboardModel {
  String get id => throw _privateConstructorUsedError;
  String get userId => throw _privateConstructorUsedError;
  String get name => throw _privateConstructorUsedError;
  Map<String, dynamic> get layout => throw _privateConstructorUsedError;
  List<WidgetModel> get widgets => throw _privateConstructorUsedError;
  bool get isPublic => throw _privateConstructorUsedError;
  DateTime get createdAt => throw _privateConstructorUsedError;
  DateTime get updatedAt => throw _privateConstructorUsedError;

  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;
  @JsonKey(ignore: true)
  $DashboardModelCopyWith<DashboardModel> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $DashboardModelCopyWith<$Res> {
  factory $DashboardModelCopyWith(
          DashboardModel value, $Res Function(DashboardModel) then) =
      _$DashboardModelCopyWithImpl<$Res, DashboardModel>;
  @useResult
  $Res call(
      {String id,
      String userId,
      String name,
      Map<String, dynamic> layout,
      List<WidgetModel> widgets,
      bool isPublic,
      DateTime createdAt,
      DateTime updatedAt});
}

/// @nodoc
class _$DashboardModelCopyWithImpl<$Res, $Val extends DashboardModel>
    implements $DashboardModelCopyWith<$Res> {
  _$DashboardModelCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? userId = null,
    Object? name = null,
    Object? layout = null,
    Object? widgets = null,
    Object? isPublic = null,
    Object? createdAt = null,
    Object? updatedAt = null,
  }) {
    return _then(_value.copyWith(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as String,
      userId: null == userId
          ? _value.userId
          : userId // ignore: cast_nullable_to_non_nullable
              as String,
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
      layout: null == layout
          ? _value.layout
          : layout // ignore: cast_nullable_to_non_nullable
              as Map<String, dynamic>,
      widgets: null == widgets
          ? _value.widgets
          : widgets // ignore: cast_nullable_to_non_nullable
              as List<WidgetModel>,
      isPublic: null == isPublic
          ? _value.isPublic
          : isPublic // ignore: cast_nullable_to_non_nullable
              as bool,
      createdAt: null == createdAt
          ? _value.createdAt
          : createdAt // ignore: cast_nullable_to_non_nullable
              as DateTime,
      updatedAt: null == updatedAt
          ? _value.updatedAt
          : updatedAt // ignore: cast_nullable_to_non_nullable
              as DateTime,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$DashboardModelImplCopyWith<$Res>
    implements $DashboardModelCopyWith<$Res> {
  factory _$$DashboardModelImplCopyWith(_$DashboardModelImpl value,
          $Res Function(_$DashboardModelImpl) then) =
      __$$DashboardModelImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {String id,
      String userId,
      String name,
      Map<String, dynamic> layout,
      List<WidgetModel> widgets,
      bool isPublic,
      DateTime createdAt,
      DateTime updatedAt});
}

/// @nodoc
class __$$DashboardModelImplCopyWithImpl<$Res>
    extends _$DashboardModelCopyWithImpl<$Res, _$DashboardModelImpl>
    implements _$$DashboardModelImplCopyWith<$Res> {
  __$$DashboardModelImplCopyWithImpl(
      _$DashboardModelImpl _value, $Res Function(_$DashboardModelImpl) _then)
      : super(_value, _then);

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? userId = null,
    Object? name = null,
    Object? layout = null,
    Object? widgets = null,
    Object? isPublic = null,
    Object? createdAt = null,
    Object? updatedAt = null,
  }) {
    return _then(_$DashboardModelImpl(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as String,
      userId: null == userId
          ? _value.userId
          : userId // ignore: cast_nullable_to_non_nullable
              as String,
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
      layout: null == layout
          ? _value._layout
          : layout // ignore: cast_nullable_to_non_nullable
              as Map<String, dynamic>,
      widgets: null == widgets
          ? _value._widgets
          : widgets // ignore: cast_nullable_to_non_nullable
              as List<WidgetModel>,
      isPublic: null == isPublic
          ? _value.isPublic
          : isPublic // ignore: cast_nullable_to_non_nullable
              as bool,
      createdAt: null == createdAt
          ? _value.createdAt
          : createdAt // ignore: cast_nullable_to_non_nullable
              as DateTime,
      updatedAt: null == updatedAt
          ? _value.updatedAt
          : updatedAt // ignore: cast_nullable_to_non_nullable
              as DateTime,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$DashboardModelImpl extends _DashboardModel {
  const _$DashboardModelImpl(
      {required this.id,
      required this.userId,
      required this.name,
      required final Map<String, dynamic> layout,
      required final List<WidgetModel> widgets,
      required this.isPublic,
      required this.createdAt,
      required this.updatedAt})
      : _layout = layout,
        _widgets = widgets,
        super._();

  factory _$DashboardModelImpl.fromJson(Map<String, dynamic> json) =>
      _$$DashboardModelImplFromJson(json);

  @override
  final String id;
  @override
  final String userId;
  @override
  final String name;
  final Map<String, dynamic> _layout;
  @override
  Map<String, dynamic> get layout {
    if (_layout is EqualUnmodifiableMapView) return _layout;
    // ignore: implicit_dynamic_type
    return EqualUnmodifiableMapView(_layout);
  }

  final List<WidgetModel> _widgets;
  @override
  List<WidgetModel> get widgets {
    if (_widgets is EqualUnmodifiableListView) return _widgets;
    // ignore: implicit_dynamic_type
    return EqualUnmodifiableListView(_widgets);
  }

  @override
  final bool isPublic;
  @override
  final DateTime createdAt;
  @override
  final DateTime updatedAt;

  @override
  String toString() {
    return 'DashboardModel(id: $id, userId: $userId, name: $name, layout: $layout, widgets: $widgets, isPublic: $isPublic, createdAt: $createdAt, updatedAt: $updatedAt)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$DashboardModelImpl &&
            (identical(other.id, id) || other.id == id) &&
            (identical(other.userId, userId) || other.userId == userId) &&
            (identical(other.name, name) || other.name == name) &&
            const DeepCollectionEquality().equals(other._layout, _layout) &&
            const DeepCollectionEquality().equals(other._widgets, _widgets) &&
            (identical(other.isPublic, isPublic) ||
                other.isPublic == isPublic) &&
            (identical(other.createdAt, createdAt) ||
                other.createdAt == createdAt) &&
            (identical(other.updatedAt, updatedAt) ||
                other.updatedAt == updatedAt));
  }

  @JsonKey(ignore: true)
  @override
  int get hashCode => Object.hash(
      runtimeType,
      id,
      userId,
      name,
      const DeepCollectionEquality().hash(_layout),
      const DeepCollectionEquality().hash(_widgets),
      isPublic,
      createdAt,
      updatedAt);

  @JsonKey(ignore: true)
  @override
  @pragma('vm:prefer-inline')
  _$$DashboardModelImplCopyWith<_$DashboardModelImpl> get copyWith =>
      __$$DashboardModelImplCopyWithImpl<_$DashboardModelImpl>(
          this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$DashboardModelImplToJson(
      this,
    );
  }
}

abstract class _DashboardModel extends DashboardModel {
  const factory _DashboardModel(
      {required final String id,
      required final String userId,
      required final String name,
      required final Map<String, dynamic> layout,
      required final List<WidgetModel> widgets,
      required final bool isPublic,
      required final DateTime createdAt,
      required final DateTime updatedAt}) = _$DashboardModelImpl;
  const _DashboardModel._() : super._();

  factory _DashboardModel.fromJson(Map<String, dynamic> json) =
      _$DashboardModelImpl.fromJson;

  @override
  String get id;
  @override
  String get userId;
  @override
  String get name;
  @override
  Map<String, dynamic> get layout;
  @override
  List<WidgetModel> get widgets;
  @override
  bool get isPublic;
  @override
  DateTime get createdAt;
  @override
  DateTime get updatedAt;
  @override
  @JsonKey(ignore: true)
  _$$DashboardModelImplCopyWith<_$DashboardModelImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

WidgetModel _$WidgetModelFromJson(Map<String, dynamic> json) {
  return _WidgetModel.fromJson(json);
}

/// @nodoc
mixin _$WidgetModel {
  String get id => throw _privateConstructorUsedError;
  String get type => throw _privateConstructorUsedError;
  Map<String, dynamic> get config => throw _privateConstructorUsedError;
  Map<String, dynamic> get position => throw _privateConstructorUsedError;

  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;
  @JsonKey(ignore: true)
  $WidgetModelCopyWith<WidgetModel> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $WidgetModelCopyWith<$Res> {
  factory $WidgetModelCopyWith(
          WidgetModel value, $Res Function(WidgetModel) then) =
      _$WidgetModelCopyWithImpl<$Res, WidgetModel>;
  @useResult
  $Res call(
      {String id,
      String type,
      Map<String, dynamic> config,
      Map<String, dynamic> position});
}

/// @nodoc
class _$WidgetModelCopyWithImpl<$Res, $Val extends WidgetModel>
    implements $WidgetModelCopyWith<$Res> {
  _$WidgetModelCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? type = null,
    Object? config = null,
    Object? position = null,
  }) {
    return _then(_value.copyWith(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as String,
      type: null == type
          ? _value.type
          : type // ignore: cast_nullable_to_non_nullable
              as String,
      config: null == config
          ? _value.config
          : config // ignore: cast_nullable_to_non_nullable
              as Map<String, dynamic>,
      position: null == position
          ? _value.position
          : position // ignore: cast_nullable_to_non_nullable
              as Map<String, dynamic>,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$WidgetModelImplCopyWith<$Res>
    implements $WidgetModelCopyWith<$Res> {
  factory _$$WidgetModelImplCopyWith(
          _$WidgetModelImpl value, $Res Function(_$WidgetModelImpl) then) =
      __$$WidgetModelImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {String id,
      String type,
      Map<String, dynamic> config,
      Map<String, dynamic> position});
}

/// @nodoc
class __$$WidgetModelImplCopyWithImpl<$Res>
    extends _$WidgetModelCopyWithImpl<$Res, _$WidgetModelImpl>
    implements _$$WidgetModelImplCopyWith<$Res> {
  __$$WidgetModelImplCopyWithImpl(
      _$WidgetModelImpl _value, $Res Function(_$WidgetModelImpl) _then)
      : super(_value, _then);

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? type = null,
    Object? config = null,
    Object? position = null,
  }) {
    return _then(_$WidgetModelImpl(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as String,
      type: null == type
          ? _value.type
          : type // ignore: cast_nullable_to_non_nullable
              as String,
      config: null == config
          ? _value._config
          : config // ignore: cast_nullable_to_non_nullable
              as Map<String, dynamic>,
      position: null == position
          ? _value._position
          : position // ignore: cast_nullable_to_non_nullable
              as Map<String, dynamic>,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$WidgetModelImpl extends _WidgetModel {
  const _$WidgetModelImpl(
      {required this.id,
      required this.type,
      required final Map<String, dynamic> config,
      required final Map<String, dynamic> position})
      : _config = config,
        _position = position,
        super._();

  factory _$WidgetModelImpl.fromJson(Map<String, dynamic> json) =>
      _$$WidgetModelImplFromJson(json);

  @override
  final String id;
  @override
  final String type;
  final Map<String, dynamic> _config;
  @override
  Map<String, dynamic> get config {
    if (_config is EqualUnmodifiableMapView) return _config;
    // ignore: implicit_dynamic_type
    return EqualUnmodifiableMapView(_config);
  }

  final Map<String, dynamic> _position;
  @override
  Map<String, dynamic> get position {
    if (_position is EqualUnmodifiableMapView) return _position;
    // ignore: implicit_dynamic_type
    return EqualUnmodifiableMapView(_position);
  }

  @override
  String toString() {
    return 'WidgetModel(id: $id, type: $type, config: $config, position: $position)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$WidgetModelImpl &&
            (identical(other.id, id) || other.id == id) &&
            (identical(other.type, type) || other.type == type) &&
            const DeepCollectionEquality().equals(other._config, _config) &&
            const DeepCollectionEquality().equals(other._position, _position));
  }

  @JsonKey(ignore: true)
  @override
  int get hashCode => Object.hash(
      runtimeType,
      id,
      type,
      const DeepCollectionEquality().hash(_config),
      const DeepCollectionEquality().hash(_position));

  @JsonKey(ignore: true)
  @override
  @pragma('vm:prefer-inline')
  _$$WidgetModelImplCopyWith<_$WidgetModelImpl> get copyWith =>
      __$$WidgetModelImplCopyWithImpl<_$WidgetModelImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$WidgetModelImplToJson(
      this,
    );
  }
}

abstract class _WidgetModel extends WidgetModel {
  const factory _WidgetModel(
      {required final String id,
      required final String type,
      required final Map<String, dynamic> config,
      required final Map<String, dynamic> position}) = _$WidgetModelImpl;
  const _WidgetModel._() : super._();

  factory _WidgetModel.fromJson(Map<String, dynamic> json) =
      _$WidgetModelImpl.fromJson;

  @override
  String get id;
  @override
  String get type;
  @override
  Map<String, dynamic> get config;
  @override
  Map<String, dynamic> get position;
  @override
  @JsonKey(ignore: true)
  _$$WidgetModelImplCopyWith<_$WidgetModelImpl> get copyWith =>
      throw _privateConstructorUsedError;
}
