import 'package:flutter/material.dart';
import '../services/api_service.dart';
import '../widgets/sos_button.dart';
import '../widgets/location_tile.dart';
import 'parking_screen.dart';
import 'sos_screen.dart';
import 'map_screen.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  final api = ApiService();
  double? lat;
  double? lng;
  String? error;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    try {
      final loc = await api.getLocation();
      setState(() {
        lat = loc['latitude'];
        lng = loc['longitude'];
      });
    } catch (e) {
      setState(() => error = e.toString());
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('TriVision Drive Safe')),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          if (lat != null && lng != null)
            Card(child: LocationTile(latitude: lat!, longitude: lng!)),
          if (error != null)
            Text(error!, style: const TextStyle(color: Colors.red)),
          const SizedBox(height: 12),
          ElevatedButton.icon(
            icon: const Icon(Icons.map),
            label: const Text('View Map'),
            onPressed: () => Navigator.push(
              context, MaterialPageRoute(builder: (_) => const MapScreen())),
          ),
          const SizedBox(height: 8),
          ElevatedButton.icon(
            icon: const Icon(Icons.local_parking),
            label: const Text('Parking'),
            onPressed: () => Navigator.push(
              context, MaterialPageRoute(builder: (_) => const ParkingScreen())),
          ),
          const SizedBox(height: 8),
          SosButton(
            onPressed: () => Navigator.push(
              context, MaterialPageRoute(builder: (_) => const SosScreen())),
          ),
        ],
      ),
    );
  }
}
