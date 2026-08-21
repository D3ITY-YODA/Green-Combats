package com.greencompass.feature.onboarding

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.MyLocation
import androidx.compose.material.icons.filled.Search
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.greencompass.core.ui.GreenCompassColors

@Composable
fun PlaceSetupScreen(onContinue: () -> Unit) {
    Surface(modifier = Modifier.fillMaxSize(), color = GreenCompassColors.WarmWhite) {
        Column(modifier = Modifier.fillMaxSize().padding(24.dp)) {
            Text(text = "Choose a place", fontSize = 24.sp, fontWeight = FontWeight.Bold, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = 8.dp))
            Text(text = "See updates for a place that matters to you.", fontSize = 15.sp, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = 32.dp))

            OutlinedButton(onClick = {}, modifier = Modifier.fillMaxWidth().height(56.dp), shape = RoundedCornerShape(12.dp)) { Icon(Icons.Default.MyLocation, contentDescription = null, tint = GreenCompassColors.ForestGreen); Spacer(Modifier.width(8.dp)); Text("Use my location", color = GreenCompassColors.ForestGreen) }
            Spacer(modifier = Modifier.height(12.dp))
            OutlinedButton(onClick = {}, modifier = Modifier.fillMaxWidth().height(56.dp), shape = RoundedCornerShape(12.dp)) { Icon(Icons.Default.Search, contentDescription = null, tint = GreenCompassColors.ForestGreen); Spacer(Modifier.width(8.dp)); Text("Search for a place", color = GreenCompassColors.ForestGreen) }
            
            Spacer(modifier = Modifier.weight(1f))
            Button(onClick = onContinue, modifier = Modifier.fillMaxWidth().height(56.dp), colors = ButtonDefaults.buttonColors(containerColor = GreenCompassColors.ForestGreen), shape = RoundedCornerShape(12.dp)) { Text("Continue", fontSize = 16.sp) }
        }
    }
}
