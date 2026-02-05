import 'package:flutter/material.dart';

class DashboardWidget extends StatelessWidget {
  final dynamic widget;
  final VoidCallback onRemove;
  final VoidCallback onConfigure;

  const DashboardWidget({
    super.key,
    required this.widget,
    required this.onRemove,
    required this.onConfigure,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Column(
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              const Text('Widget Title'),
              Row(
                children: [
                   IconButton(icon: const Icon(Icons.settings), onPressed: onConfigure),
                   IconButton(icon: const Icon(Icons.close), onPressed: onRemove),
                ],
              ),
            ],
          ),
          const Expanded(child: Center(child: Text('Content'))),
        ],
      ),
    );
  }
}
