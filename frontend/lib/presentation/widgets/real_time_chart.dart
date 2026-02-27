import 'package:flutter/material.dart';
import 'package:fl_chart/fl_chart.dart';
import 'package:financial_analytics/services/websocket_service.dart';
import 'dart:async';

class RealTimeChart extends StatefulWidget {
  final String symbol;
  final ChartType chartType;

  const RealTimeChart({
    Key? key,
    required this.symbol,
    this.chartType = ChartType.line,
  }) : super(key: key);

  @override
  State<RealTimeChart> createState() => _RealTimeChartState();
}

class _RealTimeChartState extends State<RealTimeChart> {
  final List<FlSpot> _priceData = [];
  final List<BarChartGroupData> _volumeData = [];
  StreamSubscription? _subscription;
  final _maxDataPoints = 100;
  double _minY = double.infinity;
  double _maxY = double.negativeInfinity;

  @override
  void initState() {
    super.initState();
    _subscribeToUpdates();
  }

  void _subscribeToUpdates() {
    final wsService = WebSocketService();
    wsService.subscribe([widget.symbol]);

    _subscription = wsService.messages.listen((message) {
      if (message['type'] == 'price_update' &&
          message['symbol'] == widget.symbol) {
        setState(() {
          _addDataPoint(
            message['price'].toDouble(),
            message['volume']?.toDouble() ?? 0,
            DateTime.fromMillisecondsSinceEpoch(message['timestamp']),
          );
        });
      }
    });
  }

  void _addDataPoint(double price, double volume, DateTime timestamp) {
    final x = _priceData.length.toDouble();

    _priceData.add(FlSpot(x, price));
    _volumeData.add(
      BarChartGroupData(
        x: x.toInt(),
        barRods: [
          BarChartRodData(
            toY: volume,
            color: price >
                    (_priceData.length > 1
                        ? _priceData[_priceData.length - 2].y
                        : price)
                ? Colors.green
                : Colors.red,
            width: 2,
          ),
        ],
      ),
    );

    // Update min/max for scaling
    _minY = _priceData.map((e) => e.y).reduce((a, b) => a < b ? a : b);
    _maxY = _priceData.map((e) => e.y).reduce((a, b) => a > b ? a : b);

    // Keep only recent data points
    if (_priceData.length > _maxDataPoints) {
      _priceData.removeAt(0);
      _volumeData.removeAt(0);

      // Adjust x values
      for (int i = 0; i < _priceData.length; i++) {
        _priceData[i] = FlSpot(i.toDouble(), _priceData[i].y);
      }
      for (int i = 0; i < _volumeData.length; i++) {
        _volumeData[i] = _volumeData[i].copyWith(x: i);
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  widget.symbol,
                  style: Theme.of(context).textTheme.headlineSmall,
                ),
                if (_priceData.isNotEmpty)
                  Text(
                    '\${_priceData.last.y.toStringAsFixed(2)}',
                    style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                          color: _getPriceColor(),
                        ),
                  ),
              ],
            ),
            const SizedBox(height: 16),
            Expanded(child: _buildChart()),
          ],
        ),
      ),
    );
  }

  Widget _buildChart() {
    if (_priceData.isEmpty) {
      return const Center(child: CircularProgressIndicator());
    }

    switch (widget.chartType) {
      case ChartType.line:
        return _buildLineChart();
      case ChartType.candlestick:
        return _buildCandlestickChart();
      case ChartType.combined:
        return _buildCombinedChart();
    }
  }

  Widget _buildLineChart() {
    return LineChart(
      LineChartData(
        gridData: const FlGridData(show: true),
        titlesData: const FlTitlesData(
          leftTitles: AxisTitles(
            sideTitles: SideTitles(showTitles: true, reservedSize: 40),
          ),
          bottomTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
          rightTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
          topTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
        ),
        borderData: FlBorderData(show: false),
        minY: _minY * 0.99,
        maxY: _maxY * 1.01,
        lineBarsData: [
          LineChartBarData(
            spots: _priceData,
            isCurved: true,
            color: Theme.of(context).primaryColor,
            barWidth: 2,
            isStrokeCapRound: true,
            dotData: const FlDotData(show: false),
            belowBarData: BarAreaData(
              show: true,
              color: Theme.of(context).primaryColor.withValues(alpha: 0.1),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildCandlestickChart() {
    // Implement candlestick chart
    return Container();
  }

  Widget _buildCombinedChart() {
    // Implement combined price and volume chart
    return Container();
  }

  Color _getPriceColor() {
    if (_priceData.length < 2) return Colors.grey;
    final lastPrice = _priceData.last.y;
    final previousPrice = _priceData[_priceData.length - 2].y;
    return lastPrice > previousPrice ? Colors.green : Colors.red;
  }

  @override
  void dispose() {
    _subscription?.cancel();
    WebSocketService().unsubscribe([widget.symbol]);
    super.dispose();
  }
}

enum ChartType { line, candlestick, combined }
