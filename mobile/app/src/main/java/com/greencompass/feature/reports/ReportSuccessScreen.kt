package com.greencompass.feature.reports

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.CheckCircle
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.greencompass.core.ui.GreenCompassColors

@Composable
fun ReportSuccessScreen(onDone: () -> Unit) {
    Surface(modifier = Modifier.fillMaxSize(), color = GreenCompassColors.WarmWhite) {
        Column(
            modifier = Modifier.fillMaxSize().padding(32.dp),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.Center
        ) {
            Icon(Icons.Default.CheckCircle, contentDescription = null, tint = GreenCompassColors.ForestGreen, modifier = Modifier.size(80.dp))
            Spacer(Modifier.height(24.dp))
            Text(text = "Thank you!", fontSize = 28.sp, fontWeight = FontWeight.Bold, color = GreenCompassColors.Charcoal, textAlign = TextAlign.Center)
            Spacer(Modifier.height(16.dp))
            Text(text = "Your update has been received.\nIt will be reviewed by the relevant local team.", fontSize = 16.sp, color = GreenCompassColors.MutedText, textAlign = TextAlign.Center)
            Spacer(Modifier.weight(1f))
            Button(onClick = onDone, modifier = Modifier.fillMaxWidth().height(56.dp), colors = ButtonDefaults.buttonColors(containerColor = GreenCompassColors.ForestGreen), shape = RoundedCornerShape(12.dp)) {
                Text(text = "Done", fontSize = 16.sp)
            }
        }
    }
}
