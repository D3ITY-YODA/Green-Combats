package com.greencompass.feature.profile

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowBack
import androidx.compose.material.icons.filled.ChevronRight
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.greencompass.core.ui.GreenCompassColors

@Composable
fun NotificationsScreen(onBack: () -> Unit) {
    var appNotifs by remember { mutableStateOf(true) }
    var sms by remember { mutableStateOf(false) }
    var voice by remember { mutableStateOf(false) }

    Surface(modifier = Modifier.fillMaxSize(), color = GreenCompassColors.WarmWhite) {
        Column(modifier = Modifier.fillMaxSize().padding(24.dp)) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                IconButton(onClick = onBack) { Icon(Icons.Default.ArrowBack, contentDescription = "Back", tint = GreenCompassColors.Charcoal) }
                Text(text = "Notifications", fontSize = 20.sp, fontWeight = FontWeight.SemiBold, color = GreenCompassColors.Charcoal)
            }
            Spacer(Modifier.height(24.dp))
            Text(text = "Choose how to receive updates", fontSize = 15.sp, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = 24.dp))

            NotificationRow("App notifications", appNotifs) { appNotifs = it }
            NotificationRow("SMS", sms) { sms = it }
            NotificationRow("Voice updates", voice) { voice = it }
        }
    }
}

@Composable
private fun NotificationRow(title: String, checked: Boolean, onCheckedChange: (Boolean) -> Unit) {
    Row(modifier = Modifier.fillMaxWidth().padding(vertical = 12.dp), verticalAlignment = Alignment.CenterVertically, horizontalArrangement = Arrangement.SpaceBetween) {
        Text(text = title, fontSize = 16.sp, color = GreenCompassColors.Charcoal)
        Switch(checked = checked, onCheckedChange = onCheckedChange, colors = SwitchDefaults.colors(checkedThumbColor = GreenCompassColors.ForestGreen, checkedTrackColor = GreenCompassColors.SoftSage))
    }
    Divider(color = GreenCompassColors.Stone)
}

@Composable
fun SettingsScreen(onBack: () -> Unit, onOpenNotifications: () -> Unit) {
    Surface(modifier = Modifier.fillMaxSize(), color = GreenCompassColors.WarmWhite) {
        Column(modifier = Modifier.fillMaxSize().padding(24.dp)) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                IconButton(onClick = onBack) { Icon(Icons.Default.ArrowBack, contentDescription = "Back", tint = GreenCompassColors.Charcoal) }
                Text(text = "Settings", fontSize = 20.sp, fontWeight = FontWeight.SemiBold, color = GreenCompassColors.Charcoal)
            }
            Spacer(Modifier.height(24.dp))
            
            SettingsRow("Language", "English")
            SettingsRow("Units", "Metric")
            SettingsRow("Notifications", "", onClick = onOpenNotifications)
            SettingsRow("Data usage", "Standard")
            SettingsRow("Theme", "System")
            SettingsRow("Privacy", "")
        }
    }
}

@Composable
private fun SettingsRow(title: String, subtitle: String, onClick: () -> Unit = {}) {
    Row(modifier = Modifier.fillMaxWidth().padding(vertical = 16.dp).clickable { onClick() }, verticalAlignment = Alignment.CenterVertically, horizontalArrangement = Arrangement.SpaceBetween) {
        Column {
            Text(text = title, fontSize = 16.sp, color = GreenCompassColors.Charcoal)
            if (subtitle.isNotEmpty()) Text(text = subtitle, fontSize = 13.sp, color = GreenCompassColors.MutedText)
        }
        Icon(Icons.Default.ChevronRight, contentDescription = null, tint = GreenCompassColors.MutedText)
    }
    Divider(color = GreenCompassColors.Stone)
}
