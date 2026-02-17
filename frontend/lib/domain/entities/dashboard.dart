import 'package:equatable/equatable.dart';

class Dashboard extends Equatable {
  final String id;
  final String userId;
  final String name;
  final Map<String, dynamic> layout;
  final List<DashboardWidget> widgets;
  final bool isPublic;
  final DateTime createdAt;
  final DateTime updatedAt;

  const Dashboard({
    required this.id,
    required this.userId,
    required this.name,
    required this.layout,
    required this.widgets,
    required this.isPublic,
    required this.createdAt,
    required this.updatedAt,
  });

  @override
  List<Object?> get props => [
        id,
        userId,
        name,
        layout,
        widgets,
        isPublic,
        createdAt,
        updatedAt,
      ];
}

class DashboardWidget extends Equatable {
  final String id;
  final String type;
  final Map<String, dynamic> config;
  final Map<String, dynamic> position;

  const DashboardWidget({
    required this.id,
    required this.type,
    required this.config,
    required this.position,
  });

  @override
  List<Object?> get props => [id, type, config, position];
}
