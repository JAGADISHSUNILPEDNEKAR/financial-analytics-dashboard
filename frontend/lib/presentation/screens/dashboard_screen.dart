import 'package:flutter/material.dart';
import 'package:auto_route/auto_route.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:financial_analytics/presentation/blocs/dashboard/dashboard_bloc.dart';
import 'package:financial_analytics/presentation/widgets/dashboard_widget.dart';
import 'package:financial_analytics/presentation/widgets/loading_indicator.dart';
import 'package:financial_analytics/presentation/widgets/error_widget.dart';

@RoutePage()
class DashboardScreen extends StatelessWidget {
  const DashboardScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Dashboard'),
        actions: [
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: () => _showAddWidgetDialog(context),
          ),
          IconButton(
            icon: const Icon(Icons.settings),
            onPressed: () => _navigateToSettings(context),
          ),
        ],
      ),
      body: BlocBuilder<DashboardBloc, DashboardState>(
        builder: (context, state) {
          return state.when(
            initial: () => const LoadingIndicator(),
            loading: () => const LoadingIndicator(),
            loaded: (dashboard) => _buildDashboard(context, dashboard),
            error: (message) => CustomErrorWidget(
              message: message,
              onRetry: () => context.read<DashboardBloc>().add(
                const DashboardLoadRequested(),
              ),
            ),
          );
        },
      ),
    );
  }
  
  Widget _buildDashboard(BuildContext context, Dashboard dashboard) {
    return LayoutBuilder(
      builder: (context, constraints) {
        final isDesktop = constraints.maxWidth > 1200;
        final isTablet = constraints.maxWidth > 600;
        
        return GridView.builder(
          padding: const EdgeInsets.all(16),
          gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
            crossAxisCount: isDesktop ? 4 : (isTablet ? 2 : 1),
            crossAxisSpacing: 16,
            mainAxisSpacing: 16,
            childAspectRatio: 1.5,
          ),
          itemCount: dashboard.widgets.length,
          itemBuilder: (context, index) {
            final widget = dashboard.widgets[index];
            return DashboardWidget(
              key: ValueKey(widget.id),
              widget: widget,
              onRemove: () => _removeWidget(context, widget.id),
              onConfigure: () => _configureWidget(context, widget),
            );
          },
        );
      },
    );
  }
  
  void _showAddWidgetDialog(BuildContext context) {
    // Implementation for adding widgets
  }
  
  void _navigateToSettings(BuildContext context) {
    // Navigate to settings
  }
  
  void _removeWidget(BuildContext context, String widgetId) {
    context.read<DashboardBloc>().add(
      DashboardWidgetRemoved(widgetId: widgetId),
    );
  }
  
  void _configureWidget(BuildContext context, Widget widget) {
    // Configure widget dialog
  }
}