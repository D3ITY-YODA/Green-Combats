package com.greencompass.feature.onboarding

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.greencompass.core.ui.GreenCompassColors

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun InterestsScreen(onContinue: () -> Unit, onSkip: () -> Unit) {
    var selected by remember { mutableStateOf(setOf("Local outlook", "Water outlook")) }
    val interests = listOf("Local outlook", "Water outlook", "Seasonal information", "Community updates")

    Surface(modifier = Modifier.fillMaxSize(), color = GreenCompassColors.WarmWhite) {
        Column(modifier = Modifier.fillMaxSize().padding(24.dp)) {
            Text(text = "Explore what matters to you", fontSize = 24.sp, fontWeight = FontWeight.Bold, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = 8.dp))
            Text(text = "Choose a few areas to personalize your experience.", fontSize = 15.sp, color = GreenCompassColors.MutedText, modifier = Modifier.padding(bottom = 32.dp))

            interests.forEach { interest ->
                val isSelected = selected.contains(interest)
                FilterChip(
                    selected = isSelected, 
                    onClick = { if(isSelected) selected -= interest else selected += interest }, 
                    label = { Text(interest) }, 
                    modifier = Modifier.padding(bottom = 12.dp)
                )
            }

            Spacer(modifier = Modifier.weight(1f))
            Button(onClick = onContinue, modifier = Modifier.fillMaxWidth().height(56.dp), colors = ButtonDefaults.buttonColors(containerColor = GreenCompassColors.ForestGreen), shape = RoundedCornerShape(12.dp)) { Text("Continue", fontSize = 16.sp) }
            Spacer(modifier = Modifier.height(8.dp))
            TextButton(onClick = onSkip, modifier = Modifier.fillMaxWidth()) { Text("Skip for now", color = GreenCompassColors.MutedText) }
        }
    }
}
