import 'package:flutter/material.dart';

class SosButton extends StatelessWidget {
  final VoidCallback onPressed;
  const SosButton({super.key, required this.onPressed});

  @override
  Widget build(BuildContext context) {
    return ElevatedButton(
      style: ElevatedButton.styleFrom(backgroundColor: Colors.red),
      onPressed: onPressed,
      child: const Text('SOS'),
    );
  }
}
