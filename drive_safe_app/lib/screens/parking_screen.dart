import 'package:flutter/material.dart';
import '../services/api_service.dart';
import '../services/location_service.dart';
import '../models/parking_session.dart';
import '../widgets/parking_card.dart';

class ParkingScreen extends StatefulWidget {
  const ParkingScreen({super.key});

  @override
  State<ParkingScreen> createState() => _ParkingScreenState();
}

class _ParkingScreenState extends State<ParkingScreen> {
  final api = ApiService();
  final loc = LocationService();
  final carIdController = TextEditingController(text: 'CAR123');

  List<Map<String, dynamic>> spots = [];
  ParkingSession? currentSession;
  bool loading = true;
  String? error;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    setState(() {
      loading = true;
      error = null;
    });
    try {
      final pos = await loc.currentPosition();
      final result = await api.getNearestParking(pos.latitude, pos.longitude);
      setState(() => spots = result);
    } catch (e) {
      setState(() => error = e.toString());
    } finally {
      setState(() => loading = false);
    }
  }

  Future<void> _start() async {
    try {
      final s = await api.startParking(carIdController.text.trim());
      setState(() => currentSession = s);
      if (mounted) {
        ScaffoldMessenger.of(context)
            .showSnackBar(const SnackBar(content: Text('Parking started')));
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context)
            .showSnackBar(SnackBar(content: Text('Error: $e')));
      }
    }
  }

  Future<void> _stop() async {
    final s = currentSession;
    if (s == null) return;
    try {
      final ended = await api.stopParking(s.sessionId);
      setState(() => currentSession = ended);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(
            content: Text('Stopped • ${ended.durationMinutes ?? 0} min')));
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context)
            .showSnackBar(SnackBar(content: Text('Error: $e')));
      }
    }
  }

  @override
  void dispose() {
    carIdController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Parking')),
      body: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(16),
            child: Row(
              children: [
                Expanded(
                  child: TextField(
                    controller: carIdController,
                    decoration: const InputDecoration(
                      labelText: 'Car ID',
                      border: OutlineInputBorder(),
                    ),
                  ),
                ),
                const SizedBox(width: 8),
                ElevatedButton(onPressed: _start, child: const Text('Start')),
                const SizedBox(width: 8),
                ElevatedButton(onPressed: _stop, child: const Text('Stop')),
              ],
            ),
          ),
          if (currentSession != null)
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              child: Align(
                alignment: Alignment.centerLeft,
                child: Text(
                  'Session: ${currentSession!.sessionId} • Started: ${currentSession!.startedAt}',
                  style: Theme.of(context).textTheme.bodySmall,
                ),
              ),
            ),
          const Divider(),
          if (loading)
            const Padding(
              padding: EdgeInsets.all(24),
              child: CircularProgressIndicator(),
            ),
          if (!loading && error != null)
            Padding(
              padding: const EdgeInsets.all(16),
              child: Text(error!, style: const TextStyle(color: Colors.red)),
            ),
          if (!loading && error == null)
            Expanded(
              child: ListView.builder(
                itemCount: spots.length,
                itemBuilder: (_, i) => ParkingCard(spot: spots[i]),
              ),
            ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _load,
        child: const Icon(Icons.refresh),
      ),
    );
  }
}
