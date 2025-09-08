import 'package:flutter/material.dart';
import '../services/api_service.dart';
import '../services/location_service.dart';
import '../models/sos.dart';

class SosScreen extends StatefulWidget {
  const SosScreen({super.key});

  @override
  State<SosScreen> createState() => _SosScreenState();
}

class _SosScreenState extends State<SosScreen> {
  final api = ApiService();
  final loc = LocationService();

  SosEvent? last;
  bool sending = false;
  String? error;

  Future<void> _send() async {
    setState(() {
      sending = true;
      error = null;
    });
    try {
      final pos = await loc.currentPosition();
      final res = await api.triggerSos('CAR123', pos.latitude, pos.longitude);
      setState(() => last = res);
      if (mounted) {
        ScaffoldMessenger.of(context)
            .showSnackBar(const SnackBar(content: Text('SOS sent')));
      }
    } catch (e) {
      setState(() => error = e.toString());
    } finally {
      setState(() => sending = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Emergency SOS')),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            ElevatedButton(
              style: ElevatedButton.styleFrom(backgroundColor: Colors.red),
              onPressed: sending ? null : _send,
              child: sending
                  ? const SizedBox(
                      width: 20, height: 20, child: CircularProgressIndicator(strokeWidth: 2))
                  : const Text('Send SOS'),
            ),
            const SizedBox(height: 16),
            if (error != null)
              Text(error!, style: const TextStyle(color: Colors.red)),
            if (last != null)
              Card(
                child: ListTile(
                  title: Text('Car ${last!.carId}'),
                  subtitle: Text(
                    'Lat: ${last!.latitude}, Lng: ${last!.longitude}\n'
                    'Time: ${last!.timestamp.toLocal()}',
                  ),
                ),
              ),
          ],
        ),
      ),
    );
  }
}
