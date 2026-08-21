package com.greencompass.feature.onboarding

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowBack
import androidx.compose.material.icons.filled.Check
import androidx.compose.material.icons.filled.Search
import androidx.compose.material.icons.outlined.CloudOff
import androidx.compose.material.icons.outlined.MyLocation
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.greencompass.core.ui.GreenCompassColors

@Composable
fun LocationSwitcherScreen(onBack: () -> Unit, onSelectPlace: (String) -> Unit) {
    val places = listOf("Lower Valley", "Riverside", "Upper Highland", "Lakeview", "Coastal Plains")
    var selected by remember { mutableStateOf("Lower Valley") }

    Surface(modifier = Modifier.fillMaxSize(), color = GreenCompassColors.WarmWhite) {
        Column(modifier = Modifier.fillMaxSize().padding(24.dp)) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                IconButton(onClick = onBack) { Icon(Icons.Default.ArrowBack, contentDescription = "Back", tint = GreenCompassColors.Charcoal) }
                Text(text = "Select a place", fontSize = 20.sp, fontWeight = FontWeight.SemiBold, color = GreenCompassColors.Charcoal)
            }
            Spacer(Modifier.height(24.dp))
            OutlinedTextField(value = "", onValueChange = {}, placeholder = { Text("Search places") }, leadingIcon = { Icon(Icons.Default.Search, contentDescription = null) }, modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp))
            Spacer(Modifier.height(24.dp))
            Text(text = "My places", fontSize = 14.sp, fontWeight = FontWeight.Medium, color = GreenCompassColors.MutedText)
            Spacer(Modifier.height(12.dp))
            LazyColumn(verticalArrangement = Arrangement.spacedBy(8.dp)) {
                items(places) { place ->
                    Row(modifier = Modifier.fillMaxWidth().padding(vertical = 8.dp), verticalAlignment = Alignment.CenterVertically, horizontalArrangement = Arrangement.SpaceBetween) {
                        Text(text = place, fontSize = 16.sp, color = GreenCompassColors.Charcoal)
                        if (place == selected) Icon(Icons.Default.Check, contentDescription = "Selected", tint = GreenCompassColors.ForestGreen)
                    }
                }
                item {
                    Spacer(Modifier.height(16.dp))
                    TextButton(onClick = {}, modifier = Modifier.fillMaxWidth()) { Text("+ Add new place", color = GreenCompassColors.ForestGreen) }
                }
            }
        }
    }
}

@Composable
fun PermissionScreen(onAllow: () -> Unit, onDeny: () -> Unit) {
    Surface(modifier = Modifier.fillMaxSize(), color = GreenCompassColors.WarmWhite) {
        Column(modifier = Modifier.fillMaxSize().padding(32.dp), horizontalAlignment = Alignment.CenterHorizontally, verticalArrangement = Arrangement.Center) {
            Icon(Icons.Outlined.MyLocation, contentDescription = null, tint = GreenCompassColors.ForestGreen, modifier = Modifier.size(80.dp))
            Spacer(Modifier.height(24.dp))
            Text(text = "Allow location access", fontSize = 24.sp, fontWeight = FontWeight.Bold, color = GreenCompassColors.Charcoal, textAlign = TextAlign.Center)
            Spacer(Modifier.height(16.dp))
            Text(text = "This helps us show updates relevant to your area.", fontSize = 16.sp, color = GreenCompassColors.MutedText, textAlign = TextAlign.Center)
            Spacer(Modifier.weight(1f))
            Button(onClick = onAllow, modifier = Modifier.fillMaxWidth().height(56.dp), colors = ButtonDefaults.buttonColors(containerColor = GreenCompassColors.ForestGreen), shape = RoundedCornerShape(12.dp)) { Text("Allow while using app", fontSize = 16.sp) }
            Spacer(Modifier.height(12.dp))
            TextButton(onClick = onDeny, modifier = Modifier.fillMaxWidth()) { Text("Not now", color = GreenCompassColors.MutedText) }
        }
    }
}

@Composable
fun OfflineScreen(onRetry: () -> Unit) {
    Surface(modifier = Modifier.fillMaxSize(), color = GreenCompassColors.WarmWhite) {
        Column(modifier = Modifier.fillMaxSize().padding(32.dp), horizontalAlignment = Alignment.CenterHorizontally, verticalArrangement = Arrangement.Center) {
            Icon(Icons.Outlined.CloudOff, contentDescription = null, tint = GreenCompassColors.MutedText, modifier = Modifier.size(80.dp))
            Spacer(Modifier.height(24.dp))
            Text(text = "You're offline", fontSize = 24.sp, fontWeight = FontWeight.Bold, color = GreenCompassColors.Charcoal, textAlign = TextAlign.Center)
            Spacer(Modifier.height(16.dp))
            Text(text = "Some content may be unavailable.\nWe'll show your last updated information.", fontSize = 16.sp, color = GreenCompassColors.MutedText, textAlign = TextAlign.Center)
            Spacer(Modifier.weight(1f))
            Button(onClick = onRetry, modifier = Modifier.fillMaxWidth().height(56.dp), colors = ButtonDefaults.buttonColors(containerColor = GreenCompassColors.ForestGreen), shape = RoundedCornerShape(12.dp)) { Text("Try again", fontSize = 16.sp) }
        }
    }
}
