package com.greencompass.feature.profile

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowBack
import androidx.compose.material.icons.outlined.*
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.greencompass.core.ui.GreenCompassColors

@Composable
fun HelpScreen(onBack: () -> Unit) {
    val items = listOf(
        "FAQs" to "Common questions and answers",
        "Contact support" to "We're here to help",
        "Send feedback" to "Help us improve",
        "Emergency information" to "Important contacts and resources"
    )
    Surface(modifier = Modifier.fillMaxSize(), color = GreenCompassColors.WarmWhite) {
        Column(modifier = Modifier.fillMaxSize().padding(24.dp)) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                IconButton(onClick = onBack) { Icon(Icons.Default.ArrowBack, contentDescription = "Back", tint = GreenCompassColors.Charcoal) }
                Text(text = "Help & support", fontSize = 20.sp, fontWeight = FontWeight.SemiBold, color = GreenCompassColors.Charcoal)
            }
            Spacer(Modifier.height(24.dp))
            LazyColumn(verticalArrangement = Arrangement.spacedBy(16.dp)) {
                items(items.size) { index ->
                    val (title, sub) = items[index]
                    Column(modifier = Modifier.fillMaxWidth().clickable {}) {
                        Text(text = title, fontSize = 16.sp, fontWeight = FontWeight.Medium, color = GreenCompassColors.Charcoal)
                        Text(text = sub, fontSize = 14.sp, color = GreenCompassColors.MutedText)
                        if (index < items.size - 1) Divider(modifier = Modifier.padding(top = 12.dp), color = GreenCompassColors.Stone)
                    }
                }
            }
        }
    }
}

@Composable
fun AboutScreen(onBack: () -> Unit) {
    Surface(modifier = Modifier.fillMaxSize(), color = GreenCompassColors.WarmWhite) {
        Column(modifier = Modifier.fillMaxSize().padding(32.dp), horizontalAlignment = Alignment.CenterHorizontally) {
            IconButton(onClick = onBack, modifier = Modifier.align(Alignment.Start).padding(start = 0.dp)) { Icon(Icons.Default.ArrowBack, contentDescription = "Back", tint = GreenCompassColors.Charcoal) }
            Spacer(Modifier.height(40.dp))
            Text(text = "🌿", fontSize = 48.sp)
            Spacer(Modifier.height(16.dp))
            Text(text = "GREEN COMPASS", fontSize = 18.sp, fontWeight = FontWeight.Bold, color = GreenCompassColors.DeepForest, letterSpacing = 1.sp)
            Spacer(Modifier.height(8.dp))
            Text(text = "Know Your Place. Move with Change.", fontSize = 14.sp, color = GreenCompassColors.MutedText, textAlign = TextAlign.Center)
            Spacer(Modifier.height(24.dp))
            Text(text = "Version 1.0.0", fontSize = 14.sp, color = GreenCompassColors.MutedText)
            Spacer(Modifier.height(40.dp))
            Text(text = "Green Compass provides clear, location-specific environmental updates to help communities and institutions make informed decisions.", fontSize = 15.sp, color = GreenCompassColors.Charcoal, textAlign = TextAlign.Center, lineHeight = 22.sp)
            Spacer(Modifier.weight(1f))
            Text(text = "Built with ❤ for people and the environment.", fontSize = 13.sp, color = GreenCompassColors.MutedText)
        }
    }
}
