package com.daleego.hearth.ui

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Button
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedButton
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp

@Composable
fun HomeScreen(
    deviceName: String?,
    checkInStatus: String?,
    isCheckInLoading: Boolean,
    isTracking: Boolean,
    onStartTracking: () -> Unit,
    onStopTracking: () -> Unit,
    onSendTestCheckIn: () -> Unit,
    onUnpair: () -> Unit
) {
    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(24.dp),
        verticalArrangement = Arrangement.Center,
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Text(
            text = "Hearth",
            style = MaterialTheme.typography.headlineMedium
        )

        Spacer(modifier = Modifier.height(8.dp))

        Text(
            text = "This device is paired.",
            style = MaterialTheme.typography.bodyLarge
        )

        if (!deviceName.isNullOrBlank()) {
            Spacer(modifier = Modifier.height(8.dp))
            Text(
                text = deviceName,
                style = MaterialTheme.typography.bodyMedium
            )
        }

        Spacer(modifier = Modifier.height(32.dp))

        if (isTracking) {
            Button(onClick = onStopTracking) {
                Text("Stop location tracking")
            }
        } else {
            Button(onClick = onStartTracking) {
                Text("Start location tracking")
            }
        }

        Spacer(modifier = Modifier.height(32.dp))

        Button(
            onClick = onSendTestCheckIn,
            enabled = !isCheckInLoading,
            modifier = Modifier.fillMaxWidth()
        ) {
            Text(if (isCheckInLoading) "Sending..." else "Send Test Check-In")
        }

        if (!checkInStatus.isNullOrBlank()) {
            Spacer(modifier = Modifier.height(12.dp))
            Text(
                text = checkInStatus,
                style = MaterialTheme.typography.bodyMedium
            )
        }

        Spacer(modifier = Modifier.height(12.dp))

        OutlinedButton(
            onClick = onUnpair,
            modifier = Modifier.fillMaxWidth()
        ) {
            Text("Unpair Device")
        }
    }
}
